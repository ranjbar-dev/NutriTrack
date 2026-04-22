package labresult

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
	labRepo "github.com/ranjbar-dev/nutritrack/internal/domain/labresult/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// LabResultService provides business logic for lab result management.
type LabResultService struct {
	repo     labRepo.LabResultRepository
	userRepo userRepo.UserRepository
	storage  shared.LabResultStorage
}

// NewLabResultService creates a new LabResultService.
func NewLabResultService(
	repo labRepo.LabResultRepository,
	userRepo userRepo.UserRepository,
	storage shared.LabResultStorage,
) *LabResultService {
	return &LabResultService{repo: repo, userRepo: userRepo, storage: storage}
}

// allowedMIME maps MIME types to file extensions.
var allowedMIME = map[string]string{
	"application/pdf": "pdf",
	"image/jpeg":      "jpg",
	"image/png":       "png",
}

// hasMagic checks if data starts with the given magic bytes.
func hasMagic(data []byte, magic []byte) bool {
	if len(data) < len(magic) {
		return false
	}
	return bytes.Equal(data[:len(magic)], magic)
}

// detectMIME reads the first 8 bytes to determine MIME type by magic bytes.
// Returns MIME type string and the full reader (bytes are prepended back).
func detectMIME(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 8)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", nil, err
	}
	buf = buf[:n]

	var mime string
	switch {
	case hasMagic(buf, []byte{0x25, 0x50, 0x44, 0x46}): // %PDF
		mime = "application/pdf"
	case hasMagic(buf, []byte{0xFF, 0xD8, 0xFF}): // JPEG
		mime = "image/jpeg"
	case hasMagic(buf, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}): // PNG
		mime = "image/png"
	default:
		return "", nil, shared.ErrInvalidFileType
	}

	// Prepend the consumed bytes back to the reader
	full := io.MultiReader(bytes.NewReader(buf), r)
	return mime, full, nil
}

// SubmitLabResultRequest holds parameters for submitting a lab result via file or link.
type SubmitLabResultRequest struct {
	Title      string
	ResultType string // blood_test, urine_test, thyroid, hormone, allergy, other
	TestDate   *time.Time
	Notes      string
	// File upload path (optional — one of File or Link must be provided)
	File         io.Reader
	FileOrigName string
	FileSize     int64
	// Link path (optional — takes effect when File is nil)
	Link *string
}

// SubmitLabResult validates and stores a lab result (file or link) for a client.
func (s *LabResultService) SubmitLabResult(
	ctx context.Context,
	clientID uuid.UUID,
	callerID uuid.UUID,
	callerRole string,
	req SubmitLabResultRequest,
) (*entity.LabResult, error) {
	if req.File == nil && req.Link == nil {
		return nil, shared.ErrValidation
	}

	// Load client to verify access and get NutritionistID
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, shared.ErrUserNotFound
	}

	// Determine nutritionistID and verify access
	var nutritionistID uuid.UUID
	switch callerRole {
	case "superadmin":
		if client.GetNutritionistID() == nil {
			return nil, shared.ErrInternal
		}
		nutritionistID = *client.GetNutritionistID()
	case "nutritionist":
		if !client.BelongsTo(callerID) {
			return nil, shared.ErrForbidden
		}
		nutritionistID = callerID
	case "client":
		if callerID != clientID {
			return nil, shared.ErrForbidden
		}
		if client.GetNutritionistID() == nil {
			return nil, shared.ErrInternal
		}
		nutritionistID = *client.GetNutritionistID()
	default:
		return nil, shared.ErrForbidden
	}

	result := entity.NewLabResult(clientID, nutritionistID, req.Title, req.ResultType, req.TestDate, req.Notes, req.Link)

	if req.File != nil {
		const maxSize = 10 * 1024 * 1024 // 10 MB
		if req.FileSize > maxSize {
			return nil, shared.ErrFileTooLarge
		}

		// Detect MIME type via magic bytes; get unified reader
		mimeType, fullReader, err := detectMIME(req.File)
		if err != nil {
			return nil, err
		}

		ext, ok := allowedMIME[mimeType]
		if !ok {
			return nil, shared.ErrInvalidFileType
		}

		// Save file to local storage
		filePath, err := s.storage.SaveLabResult(fullReader, ext)
		if err != nil {
			return nil, err
		}

		result.SetFilePath(filePath)
		result.SetOriginalName(req.FileOrigName)
		result.SetFileType(mimeType)
		result.SetFileSize(req.FileSize)
	}

	return s.repo.Create(ctx, result)
}

// UploadLabResult is kept for backward compatibility; delegates to SubmitLabResult.
func (s *LabResultService) UploadLabResult(
	ctx context.Context,
	clientID uuid.UUID,
	callerID uuid.UUID,
	callerRole string,
	src io.Reader,
	originalName string,
	fileSize int64,
) (*entity.LabResult, error) {
	return s.SubmitLabResult(ctx, clientID, callerID, callerRole, SubmitLabResultRequest{
		Title:        originalName,
		ResultType:   "other",
		File:         src,
		FileOrigName: originalName,
		FileSize:     fileSize,
	})
}

// ListClientLabResults returns paginated lab results for a client.
// Access control: client can view own; nutritionist can view their clients'; superadmin sees all.
func (s *LabResultService) ListClientLabResults(
	ctx context.Context,
	clientID uuid.UUID,
	callerID uuid.UUID,
	callerRole string,
	limit, offset int32,
) ([]*entity.LabResult, int64, error) {
	if err := s.checkAccess(ctx, clientID, callerID, callerRole); err != nil {
		return nil, 0, err
	}
	return s.repo.ListByClientID(ctx, clientID, limit, offset)
}

// GetLabResultForDownload retrieves a lab result for download, verifying access.
func (s *LabResultService) GetLabResultForDownload(
	ctx context.Context,
	labResultID uuid.UUID,
	callerID uuid.UUID,
	callerRole string,
) (*entity.LabResult, error) {
	result, err := s.repo.FindByID(ctx, labResultID)
	if err != nil {
		return nil, err
	}
	if err := s.checkAccess(ctx, result.ClientID(), callerID, callerRole); err != nil {
		return nil, err
	}
	return result, nil
}

// checkAccess is the shared access-control helper.
func (s *LabResultService) checkAccess(ctx context.Context, clientID, callerID uuid.UUID, callerRole string) error {
	if callerRole == "superadmin" {
		return nil
	}
	if callerRole == "client" {
		if callerID != clientID {
			return shared.ErrForbidden
		}
		return nil
	}
	// Nutritionist: verify client belongs to them
	client, err := s.userRepo.FindByID(ctx, clientID)
	if err != nil {
		return err
	}
	if client == nil {
		return shared.ErrUserNotFound
	}
	if !client.BelongsTo(callerID) {
		return shared.ErrForbidden
	}
	return nil
}
