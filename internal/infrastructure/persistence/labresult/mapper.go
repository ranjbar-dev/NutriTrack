package labresult

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.LabResult) *entity.LabResult {
	return entity.ReconstituteLabResult(
		row.ID,
		row.ClientID,
		row.NutritionistID,
		row.Title,
		row.ResultType,
		row.TestDate,
		row.FilePath,
		row.OriginalName,
		row.FileType,
		row.FileSize,
		row.Link,
		row.Notes,
		row.CreatedAt,
	)
}

// isNotFound returns true when a pgx query returns no rows.
func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
