package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
)

var (
	ErrTrackingNotFound     = errors.New("رکورد یافت نشد")
	ErrTrackingUnauthorized = errors.New("دسترسی غیرمجاز")
	ErrLabFileMissing       = errors.New("فایل یا لینک الزامی است")
	ErrLabFileInvalidType   = errors.New("فرمت فایل مجاز نیست — فقط PDF، JPG و PNG")
	ErrLabFileTooLarge      = errors.New("حجم فایل بیش از ۱۰ مگابایت مجاز است")
)

type TrackingService struct {
	repo       repository.TrackingRepository
	uploadsDir string
	logger     zerolog.Logger
}

func NewTrackingService(repo repository.TrackingRepository, uploadsDir string, logger zerolog.Logger) *TrackingService {
	return &TrackingService{repo: repo, uploadsDir: uploadsDir, logger: logger}
}

func (s *TrackingService) LogFood(ctx context.Context, clientID uuid.UUID, req dto.LogFoodRequest) (*dto.FoodLogResponse, error) {
	resp, err := s.repo.LogFood(ctx, clientID, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListFoodLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.FoodLogResponse, error) {
	resp, err := s.repo.ListFoodLogs(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListFoodLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.FoodLogResponse, error) {
	resp, err := s.repo.ListFoodLogsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) LogWater(ctx context.Context, clientID uuid.UUID, req dto.LogWaterRequest) (*dto.WaterLogResponse, error) {
	resp, err := s.repo.LogWater(ctx, clientID, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListWaterLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.WaterLogResponse, error) {
	resp, err := s.repo.ListWaterLogs(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListWaterLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WaterLogResponse, error) {
	resp, err := s.repo.ListWaterLogsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) UpsertSleep(ctx context.Context, clientID uuid.UUID, req dto.UpsertSleepRequest) (*dto.SleepLogResponse, error) {
	resp, err := s.repo.UpsertSleep(ctx, clientID, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetSleepLog(ctx context.Context, clientID uuid.UUID, date string) (*dto.SleepLogResponse, error) {
	resp, err := s.repo.GetSleepLog(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListSleepLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.SleepLogResponse, error) {
	resp, err := s.repo.ListSleepLogsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) LogExercise(ctx context.Context, clientID uuid.UUID, req dto.LogExerciseRequest) (*dto.ExerciseLogResponse, error) {
	resp, err := s.repo.LogExercise(ctx, clientID, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListExerciseLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.ExerciseLogResponse, error) {
	resp, err := s.repo.ListExerciseLogs(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListExerciseLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.ExerciseLogResponse, error) {
	resp, err := s.repo.ListExerciseLogsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) LogMedication(ctx context.Context, clientID uuid.UUID, req dto.LogMedicationRequest) (*dto.MedicationLogResponse, error) {
	resp, err := s.repo.LogMedication(ctx, clientID, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListMedicationLogs(ctx context.Context, clientID uuid.UUID, date string) ([]dto.MedicationLogResponse, error) {
	resp, err := s.repo.ListMedicationLogs(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListMedicationLogsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.MedicationLogResponse, error) {
	resp, err := s.repo.ListMedicationLogsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) UpsertBodyMeasurement(ctx context.Context, clientID, recordedBy uuid.UUID, req dto.UpsertBodyMeasurementRequest) (*dto.BodyMeasurementResponse, error) {
	resp, err := s.repo.UpsertBodyMeasurement(ctx, clientID, recordedBy, req)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetBodyMeasurement(ctx context.Context, clientID uuid.UUID, date string) (*dto.BodyMeasurementResponse, error) {
	resp, err := s.repo.GetBodyMeasurement(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListBodyMeasurements(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error) {
	resp, err := s.repo.ListBodyMeasurements(ctx, clientID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListBodyMeasurementsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.BodyMeasurementResponse, error) {
	resp, err := s.repo.ListBodyMeasurementsForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetWeightHistory(ctx context.Context, clientID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error) {
	resp, err := s.repo.GetWeightHistory(ctx, clientID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetWeightHistoryForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID, from, to string) ([]dto.WeightHistoryPointResponse, error) {
	resp, err := s.repo.GetWeightHistoryForNutritionist(ctx, clientID, nutritionistID, from, to)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListLabResults(ctx context.Context, clientID uuid.UUID) ([]dto.LabResultResponse, error) {
	resp, err := s.repo.ListLabResults(ctx, clientID)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) ListLabResultsForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID) ([]dto.LabResultResponse, error) {
	resp, err := s.repo.ListLabResultsForNutritionist(ctx, clientID, nutritionistID)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetLabResultForNutritionist(ctx context.Context, labID, clientID, nutritionistID uuid.UUID) (*dto.LabResultResponse, error) {
	resp, err := s.repo.GetLabResultForNutritionist(ctx, labID, clientID, nutritionistID)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) GetDailyDashboard(ctx context.Context, clientID uuid.UUID, date string) (*dto.DailyDashboardResponse, error) {
	resp, err := s.repo.GetDailyDashboard(ctx, clientID, date)
	return resp, s.normalizeNotFound(err)
}

func (s *TrackingService) CreateLabResult(ctx context.Context, clientID uuid.UUID, req dto.CreateLabResultRequest, fileReader io.Reader, fileSize int64, originalFilename string) (*dto.LabResultResponse, error) {
	localID, err := uuid.Parse(req.LocalID)
	if err != nil {
		return nil, fmt.Errorf("invalid local_id: %w", err)
	}

	var filePath, origFilename, mimeTypeValue *string
	var sizeBytes *int64

	if fileReader != nil {
		if fileSize > 10<<20 {
			return nil, ErrLabFileTooLarge
		}

		kind, err := mimetype.DetectReader(fileReader)
		if err != nil || !isAllowedMIME(kind.String()) {
			return nil, ErrLabFileInvalidType
		}

		if seeker, ok := fileReader.(io.Seeker); ok {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				return nil, fmt.Errorf("reset uploaded file: %w", err)
			}
		}

		ext := extensionFromMIME(kind.String())
		storedName := uuid.NewString() + ext
		relPath := filepath.Join("lab-results", clientID.String(), storedName)
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
			return nil, fmt.Errorf("save upload file: %w", err)
		}

		relValue := relPath
		nameValue := originalFilename
		mimeValue := kind.String()
		sizeValue := fileSize
		filePath = &relValue
		origFilename = &nameValue
		mimeTypeValue = &mimeValue
		sizeBytes = &sizeValue
	}

	if filePath == nil && (req.ExternalLink == nil || strings.TrimSpace(*req.ExternalLink) == "") {
		return nil, ErrLabFileMissing
	}

	resp, err := s.repo.CreateLabResult(ctx, repository.CreateLabResultParams{
		ClientID:         clientID,
		LocalID:          localID,
		UploadedBy:       clientID,
		Title:            strings.TrimSpace(req.Title),
		LabType:          req.LabType,
		TestDate:         req.TestDate,
		FilePath:         filePath,
		ExternalLink:     trimOptionalString(req.ExternalLink),
		OriginalFilename: origFilename,
		MimeType:         mimeTypeValue,
		FileSizeBytes:    sizeBytes,
	})
	if err != nil {
		return nil, s.normalizeNotFound(err)
	}
	return resp, nil
}

func (s *TrackingService) normalizeNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTrackingNotFound
	}
	return err
}

func DefaultDateRange(days int) (string, string) {
	to := time.Now().Format("2006-01-02")
	from := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	return from, to
}

func NormalizeSingleDate(date string) string {
	if strings.TrimSpace(date) != "" {
		return date
	}
	return time.Now().Format("2006-01-02")
}

func NormalizeRange(from, to string, days int) (string, string) {
	defaultFrom, defaultTo := DefaultDateRange(days)
	if strings.TrimSpace(from) == "" {
		from = defaultFrom
	}
	if strings.TrimSpace(to) == "" {
		to = defaultTo
	}
	return from, to
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func isAllowedMIME(mime string) bool {
	switch mime {
	case "application/pdf", "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}

func extensionFromMIME(mime string) string {
	switch mime {
	case "application/pdf":
		return ".pdf"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ".bin"
	}
}
