package labresult

import (
	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

func toDomain(row db.LabResult) *entity.LabResult {
	return &entity.LabResult{
		ID:             row.ID,
		ClientID:       row.ClientID,
		NutritionistID: row.NutritionistID,
		Title:          row.Title,
		ResultType:     row.ResultType,
		TestDate:       row.TestDate,
		FilePath:       row.FilePath,
		OriginalName:   row.OriginalName,
		FileType:       row.FileType,
		FileSize:       row.FileSize,
		Link:           row.Link,
		Notes:          row.Notes,
		CreatedAt:      row.CreatedAt,
	}
}

// isNotFound returns true when a pgx query returns no rows.
func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
