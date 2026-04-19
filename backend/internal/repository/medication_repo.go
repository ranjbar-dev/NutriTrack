package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// MedicationRepository defines operations on medications.
type MedicationRepository interface {
	Create(ctx context.Context, params sqlc.CreateMedicationParams) (*sqlc.Medication, error)
	GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetMedicationByIDRow, error)
	List(ctx context.Context, params sqlc.ListMedicationsParams) ([]sqlc.ListMedicationsRow, error)
	Count(ctx context.Context, params sqlc.CountMedicationsParams) (int64, error)
	Update(ctx context.Context, params sqlc.UpdateMedicationParams) (*sqlc.Medication, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SoftDeleteByOwner(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error
	CheckDuplicate(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
}

type medicationRepository struct {
	q *sqlc.Queries
}

// NewMedicationRepository creates a new MedicationRepository backed by the given sqlc.DBTX.
func NewMedicationRepository(db sqlc.DBTX) MedicationRepository {
	return &medicationRepository{q: sqlc.New(db)}
}

func (r *medicationRepository) Create(ctx context.Context, params sqlc.CreateMedicationParams) (*sqlc.Medication, error) {
	med, err := r.q.CreateMedication(ctx, params)
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *medicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetMedicationByIDRow, error) {
	med, err := r.q.GetMedicationByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *medicationRepository) List(ctx context.Context, params sqlc.ListMedicationsParams) ([]sqlc.ListMedicationsRow, error) {
	return r.q.ListMedications(ctx, params)
}

func (r *medicationRepository) Count(ctx context.Context, params sqlc.CountMedicationsParams) (int64, error) {
	return r.q.CountMedications(ctx, params)
}

func (r *medicationRepository) Update(ctx context.Context, params sqlc.UpdateMedicationParams) (*sqlc.Medication, error) {
	med, err := r.q.UpdateMedication(ctx, params)
	if err != nil {
		return nil, err
	}
	return &med, nil
}

func (r *medicationRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteMedication(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *medicationRepository) SoftDeleteByOwner(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	return r.q.SoftDeleteMedicationByOwner(ctx, sqlc.SoftDeleteMedicationByOwnerParams{
		ID:        pgtype.UUID{Bytes: id, Valid: true},
		CreatedBy: pgtype.UUID{Bytes: ownerID, Valid: true},
	})
}

func (r *medicationRepository) CheckDuplicate(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	arg := sqlc.CheckDuplicateMedicationNameParams{
		Name:      name,
		ExcludeID: pgtype.UUID{Valid: false},
	}
	if excludeID != nil {
		arg.ExcludeID = pgtype.UUID{Bytes: *excludeID, Valid: true}
	}
	return r.q.CheckDuplicateMedicationName(ctx, arg)
}
