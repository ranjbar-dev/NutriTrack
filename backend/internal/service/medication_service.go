package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

var (
	ErrMedicationDuplicate          = errors.New("دارو با این نام قبلاً ثبت شده است")
	ErrMedicationInvalidName        = errors.New("نام دارو الزامی است")
	ErrMedicationNotFound           = errors.New("دارو یافت نشد")
	ErrMedicationUnauthorizedEdit   = errors.New("شما مجوز ویرایش این دارو را ندارید")
	ErrMedicationUnauthorizedDelete = errors.New("شما مجوز حذف این دارو را ندارید")
)

// MedicationService handles medication management business logic.
type MedicationService struct {
	medRepo repository.MedicationRepository
	logger  zerolog.Logger
}

// NewMedicationService creates a new MedicationService with the given dependencies.
func NewMedicationService(medRepo repository.MedicationRepository, logger zerolog.Logger) *MedicationService {
	return &MedicationService{
		medRepo: medRepo,
		logger:  logger,
	}
}

// CreateMedication creates a new medication after validating for duplicates.
func (s *MedicationService) CreateMedication(ctx context.Context, userID uuid.UUID, req dto.CreateMedicationRequest) (*dto.MedicationResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrMedicationInvalidName
	}
	req.Name = name

	isDuplicate, err := s.medRepo.CheckDuplicate(ctx, name, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("name", name).Msg("failed to check duplicate medication name")
		return nil, fmt.Errorf("check duplicate medication name: %w", err)
	}
	if isDuplicate {
		return nil, ErrMedicationDuplicate
	}

	params := sqlc.CreateMedicationParams{
		Name:        name,
		GenericName: optionalText(req.GenericName),
		Form:        sqlc.MedicationForm(req.Form),
		DosageUnit:  optionalText(req.DosageUnit),
		Description: optionalText(req.Description),
		CreatedBy:   pgtype.UUID{Bytes: userID, Valid: true},
	}

	med, err := s.medRepo.Create(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID.String()).Msg("failed to create medication")
		return nil, fmt.Errorf("create medication: %w", err)
	}

	medID := uuid.UUID(med.ID.Bytes)
	return s.GetMedication(ctx, medID)
}

// GetMedication retrieves a medication by ID.
func (s *MedicationService) GetMedication(ctx context.Context, id uuid.UUID) (*dto.MedicationResponse, error) {
	med, err := s.medRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMedicationNotFound
		}
		s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to get medication")
		return nil, fmt.Errorf("get medication: %w", err)
	}

	return medicationRowToResponse(med), nil
}

// ListMedications returns a paginated list of medications with optional Persian search.
func (s *MedicationService) ListMedications(ctx context.Context, query dto.MedicationListQueryParams) (*dto.MedicationListResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	listParams := sqlc.ListMedicationsParams{
		IsActive:  optionalBool(query.IsActive),
		Search:    optionalText(query.Search),
		OffsetVal: int64((page - 1) * limit),
		LimitVal:  int64(limit),
	}

	countParams := sqlc.CountMedicationsParams{
		IsActive: optionalBool(query.IsActive),
		Search:   optionalText(query.Search),
	}

	meds, err := s.medRepo.List(ctx, listParams)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list medications")
		return nil, fmt.Errorf("list medications: %w", err)
	}

	total, err := s.medRepo.Count(ctx, countParams)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to count medications")
		return nil, fmt.Errorf("count medications: %w", err)
	}

	data := make([]dto.MedicationResponse, 0, len(meds))
	for _, med := range meds {
		data = append(data, *medicationListRowToResponse(med))
	}

	return &dto.MedicationListResponse{
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(page*limit) < total,
	}, nil
}

// UpdateMedication updates a medication with row-level authorization.
func (s *MedicationService) UpdateMedication(ctx context.Context, id, userID uuid.UUID, role string, req dto.UpdateMedicationRequest) (*dto.MedicationResponse, error) {
	current, err := s.medRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMedicationNotFound
		}
		s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to load medication for update")
		return nil, fmt.Errorf("load medication for update: %w", err)
	}

	if role == string(model.RoleNutritionist) && uuid.UUID(current.CreatedBy.Bytes) != userID {
		return nil, ErrMedicationUnauthorizedEdit
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrMedicationInvalidName
	}
	req.Name = name

	isDuplicate, err := s.medRepo.CheckDuplicate(ctx, name, &id)
	if err != nil {
		s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to check duplicate medication name for update")
		return nil, fmt.Errorf("check duplicate medication name: %w", err)
	}
	if isDuplicate {
		return nil, ErrMedicationDuplicate
	}

	params := sqlc.UpdateMedicationParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		Name:        name,
		GenericName: optionalText(req.GenericName),
		Form:        sqlc.MedicationForm(req.Form),
		DosageUnit:  optionalText(req.DosageUnit),
		Description: optionalText(req.Description),
	}

	if _, err := s.medRepo.Update(ctx, params); err != nil {
		s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to update medication")
		return nil, fmt.Errorf("update medication: %w", err)
	}

	return s.GetMedication(ctx, id)
}

// DeleteMedication soft-deletes a medication with row-level authorization.
func (s *MedicationService) DeleteMedication(ctx context.Context, id, userID uuid.UUID, role string) error {
	current, err := s.medRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMedicationNotFound
		}
		s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to load medication for delete")
		return fmt.Errorf("load medication for delete: %w", err)
	}

	switch role {
	case string(model.RoleSuperAdmin):
		if err := s.medRepo.SoftDelete(ctx, id); err != nil {
			s.logger.Error().Err(err).Str("medication_id", id.String()).Msg("failed to soft delete medication")
			return fmt.Errorf("soft delete medication: %w", err)
		}
	case string(model.RoleNutritionist):
		if uuid.UUID(current.CreatedBy.Bytes) != userID {
			return ErrMedicationUnauthorizedDelete
		}
		if err := s.medRepo.SoftDeleteByOwner(ctx, id, userID); err != nil {
			s.logger.Error().Err(err).Str("medication_id", id.String()).Str("user_id", userID.String()).Msg("failed to soft delete own medication")
			return fmt.Errorf("soft delete medication by owner: %w", err)
		}
	default:
		return ErrMedicationUnauthorizedDelete
	}

	s.logger.Info().
		Str("medication_id", id.String()).
		Str("deleted_by", userID.String()).
		Str("role", role).
		Str("created_by", uuid.UUID(current.CreatedBy.Bytes).String()).
		Msg("medication deleted")

	return nil
}

func medicationRowToResponse(med *sqlc.GetMedicationByIDRow) *dto.MedicationResponse {
	resp := &dto.MedicationResponse{
		ID:          uuid.UUID(med.ID.Bytes).String(),
		Name:        med.Name,
		Form:        string(med.Form),
		IsActive:    med.IsActive,
		CreatedBy:   uuid.UUID(med.CreatedBy.Bytes).String(),
		CreatorName: med.CreatorName,
		CreatedAt:   formatTimestamp(med.CreatedAt),
		UpdatedAt:   formatTimestamp(med.UpdatedAt),
	}
	if med.GenericName.Valid {
		resp.GenericName = &med.GenericName.String
	}
	if med.DosageUnit.Valid {
		resp.DosageUnit = &med.DosageUnit.String
	}
	if med.Description.Valid {
		resp.Description = &med.Description.String
	}
	return resp
}

func medicationListRowToResponse(med sqlc.ListMedicationsRow) *dto.MedicationResponse {
	resp := &dto.MedicationResponse{
		ID:          uuid.UUID(med.ID.Bytes).String(),
		Name:        med.Name,
		Form:        string(med.Form),
		IsActive:    med.IsActive,
		CreatedBy:   uuid.UUID(med.CreatedBy.Bytes).String(),
		CreatorName: med.CreatorName,
		CreatedAt:   formatTimestamp(med.CreatedAt),
		UpdatedAt:   formatTimestamp(med.UpdatedAt),
	}
	if med.GenericName.Valid {
		resp.GenericName = &med.GenericName.String
	}
	if med.DosageUnit.Valid {
		resp.DosageUnit = &med.DosageUnit.String
	}
	if med.Description.Valid {
		resp.Description = &med.Description.String
	}
	return resp
}


