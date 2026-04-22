package medication

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgMedicationRepository implements domain/medication/repository.MedicationRepository using PostgreSQL + sqlc.
type PgMedicationRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewPgMedicationRepository constructs a PgMedicationRepository.
func NewPgMedicationRepository(pool *pgxpool.Pool) *PgMedicationRepository {
	return &PgMedicationRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// Create inserts a new medication row and populates the entity with DB-generated fields.
func (r *PgMedicationRepository) Create(ctx context.Context, med *entity.Medication) error {
	created, err := r.queries.CreateMedication(ctx, db.CreateMedicationParams{
		Name:           med.Name(),
		NameNormalized: med.NameNormalized(),
		Description:    med.Description(),
		Unit:           med.Unit(),
		CreatedBy:      uuidToPgtypeUUID(med.CreatedBy()),
	})
	if err != nil {
		return shared.ErrInternal
	}

	med.SetID(created.ID)
	med.SetIsActive(created.IsActive)
	med.SetCreatedAt(created.CreatedAt)
	med.SetUpdatedAt(created.UpdatedAt)

	return nil
}

// FindByID retrieves a medication by its ID.
// Returns nil, nil when not found.
func (r *PgMedicationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Medication, error) {
	row, err := r.queries.GetMedicationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return toDomain(row), nil
}

// Update persists updated medication fields.
func (r *PgMedicationRepository) Update(ctx context.Context, med *entity.Medication) error {
	updated, err := r.queries.UpdateMedication(ctx, db.UpdateMedicationParams{
		ID:             med.ID(),
		Name:           med.Name(),
		NameNormalized: med.NameNormalized(),
		Description:    med.Description(),
		Unit:           med.Unit(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrMedicationNotFound
		}
		return shared.ErrInternal
	}
	med.SetUpdatedAt(updated.UpdatedAt)
	return nil
}

// Delete hard-deletes a medication record (superadmin only).
func (r *PgMedicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteMedication(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// Deactivate soft-deletes a medication record (nutritionist).
func (r *PgMedicationRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeactivateMedication(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

// Search returns active medications matching the query with pg_trgm similarity.
func (r *PgMedicationRepository) Search(ctx context.Context, query string, limit, offset int32) ([]*entity.Medication, error) {
	rows, err := r.queries.SearchMedications(ctx, db.SearchMedicationsParams{
		Query: query,
		Lim:   limit,
		Off:   offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomainList(rows), nil
}

// CountSearch returns the total count of active medications matching the query.
func (r *PgMedicationRepository) CountSearch(ctx context.Context, query string) (int64, error) {
	count, err := r.queries.CountSearchMedications(ctx, query)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}

// uuidToPgtypeUUID converts *uuid.UUID to pgtype.UUID.
func uuidToPgtypeUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: [16]byte(*id), Valid: true}
}
