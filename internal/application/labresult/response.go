package labresult

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
)

// MapLabResultResponse converts a LabResult entity to a JSON-serialisable map.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapLabResultResponse(r *entity.LabResult) map[string]any {
	m := map[string]any{
		"id":              r.ID(),
		"client_id":       r.ClientID(),
		"nutritionist_id": r.NutritionistID(),
		"title":           r.Title(),
		"result_type":     r.ResultType(),
		"test_date":       nil,
		"original_name":   r.OriginalName(),
		"file_type":       r.FileType(),
		"file_size":       r.FileSize(),
		"link":            r.Link(),
		"notes":           r.Notes(),
		"created_at":      r.CreatedAt(),
	}
	if r.TestDate() != nil {
		m["test_date"] = r.TestDate().Format("2006-01-02")
	}
	return m
}
