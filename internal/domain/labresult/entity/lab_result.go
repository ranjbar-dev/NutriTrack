package entity

import (
	"time"

	"github.com/google/uuid"
)

// LabResult represents a lab result file uploaded by or for a client.
type LabResult struct {
	ID             uuid.UUID
	ClientID       uuid.UUID
	NutritionistID uuid.UUID
	Title          string
	ResultType     string     // blood_test, urine_test, thyroid, hormone, allergy, other
	TestDate       *time.Time
	FilePath       string     // filesystem path; empty when link is used
	OriginalName   string     // original filename for Content-Disposition
	FileType       string     // MIME type: application/pdf, image/jpeg, image/png
	FileSize       int64      // bytes; 0 when link is used
	Link           *string    // URL link as alternative to file upload
	Notes          string
	CreatedAt      time.Time
}
