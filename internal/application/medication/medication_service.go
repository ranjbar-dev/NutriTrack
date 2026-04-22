package medication

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
	medRepo "github.com/ranjbar-dev/nutritrack/internal/domain/medication/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

type CreateMedicationRequest struct {
	Name        string
	Description string
	Unit        string
	CallerID    uuid.UUID
	CallerRole  string
}

type UpdateMedicationRequest struct {
	ID          uuid.UUID
	Name        string
	Description string
	Unit        string
	CallerID    uuid.UUID
	CallerRole  string
}

// MedicationService handles medication/supplement business logic.
type MedicationService struct {
	medRepo medRepo.MedicationRepository
}

// NewMedicationService constructs a MedicationService.
func NewMedicationService(repo medRepo.MedicationRepository) *MedicationService {
	return &MedicationService{medRepo: repo}
}

// CreateMedication creates a new medication/supplement entry.
func (s *MedicationService) CreateMedication(ctx context.Context, req CreateMedicationRequest) (*entity.Medication, error) {
	normalized := shared.NormalizePersian(req.Name)
	med := entity.NewMedication(req.Name, normalized, req.Description, req.Unit, &req.CallerID)
	if err := s.medRepo.Create(ctx, med); err != nil {
		return nil, err
	}
	return med, nil
}

// GetMedication fetches an active medication by ID.
func (s *MedicationService) GetMedication(ctx context.Context, id uuid.UUID) (*entity.Medication, error) {
	med, err := s.medRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if med == nil || !med.IsActive() {
		return nil, shared.ErrMedicationNotFound
	}
	return med, nil
}

// SearchMedications performs pg_trgm similarity search on active medications.
func (s *MedicationService) SearchMedications(ctx context.Context, query string, limit, offset int32) ([]*entity.Medication, int64, error) {
	normalized := shared.NormalizePersian(query)
	meds, err := s.medRepo.Search(ctx, normalized, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.medRepo.CountSearch(ctx, normalized)
	if err != nil {
		return nil, 0, err
	}
	return meds, total, nil
}

// UpdateMedication updates a medication.
// Nutritionist can only update medications they created; superadmin can update any.
func (s *MedicationService) UpdateMedication(ctx context.Context, req UpdateMedicationRequest) (*entity.Medication, error) {
	med, err := s.medRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if med == nil || !med.IsActive() {
		return nil, shared.ErrMedicationNotFound
	}
	if req.CallerRole == "nutritionist" {
		if med.CreatedBy() == nil || *med.CreatedBy() != req.CallerID {
			return nil, shared.ErrForbidden
		}
	} else if req.CallerRole != "superadmin" {
		return nil, shared.ErrForbidden
	}
	med.SetName(req.Name)
	med.SetNameNormalized(shared.NormalizePersian(req.Name))
	med.SetDescription(req.Description)
	med.SetUnit(req.Unit)
	if err := s.medRepo.Update(ctx, med); err != nil {
		return nil, err
	}
	return med, nil
}

// DeleteMedication hard-deletes (superadmin) or deactivates (nutritionist own) a medication.
func (s *MedicationService) DeleteMedication(ctx context.Context, id uuid.UUID, callerID uuid.UUID, callerRole string) error {
	med, err := s.medRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if med == nil {
		return shared.ErrMedicationNotFound
	}
	switch callerRole {
	case "superadmin":
		return s.medRepo.Delete(ctx, id)
	case "nutritionist":
		if med.CreatedBy() == nil || *med.CreatedBy() != callerID {
			return shared.ErrForbidden
		}
		return s.medRepo.Deactivate(ctx, id)
	default:
		return shared.ErrForbidden
	}
}
