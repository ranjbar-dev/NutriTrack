package labresult

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgLabResultRepository is the PostgreSQL implementation of LabResultRepository.
type PgLabResultRepository struct {
	queries *db.Queries
}

// NewPgLabResultRepository creates a new PgLabResultRepository.
func NewPgLabResultRepository(pool *pgxpool.Pool) *PgLabResultRepository {
	return &PgLabResultRepository{queries: db.New(pool)}
}

// Create inserts a new lab result record.
func (r *PgLabResultRepository) Create(ctx context.Context, result *entity.LabResult) (*entity.LabResult, error) {
	row, err := r.queries.CreateLabResult(ctx, db.CreateLabResultParams{
		ClientID:       result.ClientID(),
		NutritionistID: result.NutritionistID(),
		Title:          result.Title(),
		ResultType:     result.ResultType(),
		TestDate:       result.TestDate(),
		FilePath:       result.FilePath(),
		OriginalName:   result.OriginalName(),
		FileType:       result.FileType(),
		FileSize:       result.FileSize(),
		Link:           result.Link(),
		Notes:          result.Notes(),
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomain(row), nil
}

// FindByID retrieves a lab result by its ID.
func (r *PgLabResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.LabResult, error) {
	row, err := r.queries.GetLabResultByID(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return nil, shared.ErrLabResultNotFound
		}
		return nil, shared.ErrInternal
	}
	return toDomain(row), nil
}

// ListByClientID returns paginated lab results for a client along with total count.
func (r *PgLabResultRepository) ListByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.LabResult, int64, error) {
	rows, err := r.queries.ListLabResultsByClientID(ctx, db.ListLabResultsByClientIDParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	total, err := r.queries.CountLabResultsByClientID(ctx, clientID)
	if err != nil {
		return nil, 0, shared.ErrInternal
	}
	results := make([]*entity.LabResult, len(rows))
	for i, row := range rows {
		results[i] = toDomain(row)
	}
	return results, total, nil
}
