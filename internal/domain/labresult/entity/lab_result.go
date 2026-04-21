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
	FilePath       string // filesystem path to serve the file
	OriginalName   string // original filename for Content-Disposition
	FileType       string // MIME type: application/pdf, image/jpeg, image/png
	FileSize       int64  // bytes
	Notes          string
	CreatedAt      time.Time
}
