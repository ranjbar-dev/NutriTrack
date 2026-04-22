package medication

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
)

// MapMedicationResponse converts a domain Medication to a JSON-serialisable map.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapMedicationResponse(m *entity.Medication) map[string]any {
	resp := map[string]any{
		"id":          m.ID(),
		"name":        m.Name(),
		"description": m.Description(),
		"unit":        m.Unit(),
		"is_active":   m.IsActive(),
		"created_at":  m.CreatedAt(),
		"updated_at":  m.UpdatedAt(),
	}
	if m.CreatedBy() != nil {
		resp["created_by"] = *m.CreatedBy()
	} else {
		resp["created_by"] = nil
	}
	return resp
}
