package medication

import (
	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// toDomain converts a sqlc Medication row to a domain entity.
func toDomain(m db.Medication) *entity.Medication {
	var createdBy *uuid.UUID
	if m.CreatedBy.Valid {
		id := uuid.UUID(m.CreatedBy.Bytes)
		createdBy = &id
	}
	return entity.ReconstituteMedication(
		m.ID,
		m.Name,
		m.NameNormalized,
		m.Description,
		m.Unit,
		createdBy,
		m.IsActive,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// toDomainList converts a slice of sqlc Medication rows to domain entities.
func toDomainList(ms []db.Medication) []*entity.Medication {
	result := make([]*entity.Medication, len(ms))
	for i, m := range ms {
		result[i] = toDomain(m)
	}
	return result
}
