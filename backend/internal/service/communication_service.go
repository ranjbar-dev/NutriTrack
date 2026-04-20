package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
)

var (
	ErrCommNotFound            = errors.New("رکورد یافت نشد")
	ErrCommUnauthorized        = errors.New("دسترسی غیرمجاز")
	ErrMsgAttachmentTooLarge   = errors.New("حجم فایل بیش از حد مجاز است")
	ErrMsgAttachmentInvalidType = errors.New("فرمت فایل مجاز نیست — فقط JPG، PNG و PDF")
	ErrMsgNoContent            = errors.New("پیام باید متن یا فایل داشته باشد")
	ErrFoodRequestAlreadyReviewed = errors.New("این درخواست قبلاً بررسی شده است")
)

type CommunicationService struct {
	repo       repository.CommunicationRepository
	userRepo   repository.UserRepository
	uploadsDir string
	logger     zerolog.Logger
	notifSvc   NotificationService
}

func NewCommunicationService(repo repository.CommunicationRepository, userRepo repository.UserRepository, uploadsDir string, logger zerolog.Logger, notifSvc NotificationService) *CommunicationService {
	return &CommunicationService{repo: repo, userRepo: userRepo, uploadsDir: uploadsDir, logger: logger, notifSvc: notifSvc}
}

// ─── Messaging ───────────────────────────────────────────────────────────────

// SendMessageTo sends a message from senderID to receiverID after verifying they are a valid pair.
func (s *CommunicationService) SendMessageTo(ctx context.Context, senderID, receiverID uuid.UUID, content *string, fileReader io.Reader, fileSize int64, filename, mimeStr string) (*dto.MessageResponse, error) {
	if err := s.verifyConversationPair(ctx, senderID, receiverID); err != nil {
		return nil, err
	}

	if content == nil && fileReader == nil {
		return nil, ErrMsgNoContent
	}

	var attachmentType, attachmentPath, attachmentName *string

	if fileReader != nil {
		imageLimit := int64(5 << 20)
		pdfLimit := int64(10 << 20)

		kind, err := mimetype.DetectReader(fileReader)
		if err != nil || !isMsgAllowedMIME(kind.String()) {
			return nil, ErrMsgAttachmentInvalidType
		}

		isPDF := kind.Extension() == ".pdf"
		sizeLimit := imageLimit
		if isPDF {
			sizeLimit = pdfLimit
		}
		if fileSize > sizeLimit {
			return nil, ErrMsgAttachmentTooLarge
		}

		if seeker, ok := fileReader.(io.Seeker); ok {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				return nil, fmt.Errorf("reset file: %w", err)
			}
		}

		ext := kind.Extension()
		storedName := uuid.NewString() + ext
		// Store in sender's subdirectory
		relPath := filepath.Join("messages", senderID.String(), storedName)
		absPath := filepath.Join(s.uploadsDir, relPath)
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			return nil, fmt.Errorf("create upload dir: %w", err)
		}
		dst, err := os.Create(absPath)
		if err != nil {
			return nil, fmt.Errorf("create upload file: %w", err)
		}
		defer dst.Close()
		if _, err := io.Copy(dst, fileReader); err != nil {
			return nil, fmt.Errorf("save file: %w", err)
		}

		aType := "image"
		if isPDF {
			aType = "pdf"
		}
		attachmentType = &aType
		attachmentPath = &relPath
		attachmentName = &filename
	}

	msg, err := s.repo.CreateMessage(ctx, senderID, receiverID, content, attachmentType, attachmentPath, attachmentName)
	if err != nil {
		return nil, err
	}

	// D-18: fire-and-forget push notification to receiver
	go s.notifSvc.SendToClient(context.Background(), receiverID.String(), "new_message", dto.PushPayload{
		Type:  "new_message",
		Title: "پیام جدید",
		Body:  truncateMsgContent(content, 80),
		URL:   "/client/messages",
	})

	return msg, nil
}

// truncateMsgContent returns a truncated preview of message content.
func truncateMsgContent(s *string, n int) string {
	if s == nil {
		return "پیوست"
	}
	if len(*s) <= n {
		return *s
	}
	return (*s)[:n] + "..."
}

func (s *CommunicationService) GetMessages(ctx context.Context, requestorID, partnerID uuid.UUID, limit, offset int32) ([]dto.MessageResponse, error) {
	if err := s.verifyConversationPair(ctx, requestorID, partnerID); err != nil {
		return nil, err
	}
	return s.repo.ListMessages(ctx, requestorID, partnerID, limit, offset)
}

func (s *CommunicationService) GetNewMessages(ctx context.Context, requestorID, partnerID uuid.UUID, since time.Time) ([]dto.MessageResponse, error) {
	if err := s.verifyConversationPair(ctx, requestorID, partnerID); err != nil {
		return nil, err
	}
	return s.repo.ListMessagesSince(ctx, requestorID, partnerID, since)
}

func (s *CommunicationService) MarkRead(ctx context.Context, receiverID, senderID uuid.UUID) error {
	return s.repo.MarkMessagesRead(ctx, receiverID, senderID)
}

func (s *CommunicationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error) {
	count, err := s.repo.CountUnreadMessages(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.UnreadCountResponse{Count: count}, nil
}

func (s *CommunicationService) GetMessageAttachment(ctx context.Context, requestorID, messageID uuid.UUID) (string, string, error) {
	msg, err := s.repo.GetMessageByID(ctx, messageID)
	if err != nil {
		return "", "", s.normalizeNotFound(err)
	}
	senderID := uuid.UUID(msg.SenderID.Bytes)
	receiverID := uuid.UUID(msg.ReceiverID.Bytes)
	if requestorID != senderID && requestorID != receiverID {
		return "", "", ErrCommUnauthorized
	}
	if !msg.AttachmentPath.Valid || !msg.AttachmentName.Valid {
		return "", "", ErrCommNotFound
	}
	absPath := filepath.Join(s.uploadsDir, msg.AttachmentPath.String)
	return absPath, msg.AttachmentName.String, nil
}

// ─── Food Requests ────────────────────────────────────────────────────────────

func (s *CommunicationService) CreateFoodRequest(ctx context.Context, clientID uuid.UUID, req dto.FoodRequestCreateRequest) (*dto.FoodRequestResponse, error) {
	return s.repo.CreateFoodRequest(ctx, clientID, req.FoodName, req.Description)
}

func (s *CommunicationService) ListClientFoodRequests(ctx context.Context, clientID uuid.UUID) ([]dto.FoodRequestResponse, error) {
	return s.repo.ListFoodRequestsByClient(ctx, clientID)
}

func (s *CommunicationService) ListNutriPendingFoodRequests(ctx context.Context, nutritionistID uuid.UUID) ([]dto.FoodRequestResponse, error) {
	return s.repo.ListPendingFoodRequestsForNutritionist(ctx, nutritionistID)
}

func (s *CommunicationService) ApproveFoodRequest(ctx context.Context, id, nutritionistID uuid.UUID) (*dto.FoodRequestResponse, error) {
	// Verify the request belongs to a client of this nutritionist
	fr, err := s.repo.GetFoodRequestByID(ctx, id)
	if err != nil {
		return nil, s.normalizeNotFound(err)
	}
	clientID := uuid.UUID(fr.RequestedBy.Bytes)
	client, err := s.userRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, s.normalizeNotFound(err)
	}
	if !client.NutritionistID.Valid || uuid.UUID(client.NutritionistID.Bytes) != nutritionistID {
		return nil, ErrCommUnauthorized
	}
	result, err := s.repo.ApproveFoodRequest(ctx, id, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFoodRequestAlreadyReviewed
		}
		return nil, err
	}
	return result, nil
}

func (s *CommunicationService) RejectFoodRequest(ctx context.Context, id, nutritionistID uuid.UUID, reason *string) (*dto.FoodRequestResponse, error) {
	fr, err := s.repo.GetFoodRequestByID(ctx, id)
	if err != nil {
		return nil, s.normalizeNotFound(err)
	}
	clientID := uuid.UUID(fr.RequestedBy.Bytes)
	client, err := s.userRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, s.normalizeNotFound(err)
	}
	if !client.NutritionistID.Valid || uuid.UUID(client.NutritionistID.Bytes) != nutritionistID {
		return nil, ErrCommUnauthorized
	}
	result, err := s.repo.RejectFoodRequest(ctx, id, reason, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFoodRequestAlreadyReviewed
		}
		return nil, err
	}
	return result, nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// verifyConversationPair checks that userA and userB are a valid nutritionist-client pair.
func (s *CommunicationService) verifyConversationPair(ctx context.Context, userA, userB uuid.UUID) error {
	a, err := s.userRepo.GetByID(ctx, userA)
	if err != nil {
		return s.normalizeNotFound(err)
	}
	b, err := s.userRepo.GetByID(ctx, userB)
	if err != nil {
		return s.normalizeNotFound(err)
	}

	aIsClient := string(a.Role) == "client"
	bIsClient := string(b.Role) == "client"

	switch {
	case aIsClient && !bIsClient:
		// a is client, b should be a's nutritionist
		if !a.NutritionistID.Valid || uuid.UUID(a.NutritionistID.Bytes) != userB {
			return ErrCommUnauthorized
		}
	case !aIsClient && bIsClient:
		// a is nutritionist, b should be a's client
		if !b.NutritionistID.Valid || uuid.UUID(b.NutritionistID.Bytes) != userA {
			return ErrCommUnauthorized
		}
	default:
		// same role — not allowed
		return ErrCommUnauthorized
	}
	return nil
}

func (s *CommunicationService) normalizeNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCommNotFound
	}
	return err
}

func isMsgAllowedMIME(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png", "application/pdf":
		return true
	}
	return false
}
