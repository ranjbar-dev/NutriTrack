package entity

import (
	"time"

	"github.com/google/uuid"
)

// LabResult represents a lab result file uploaded by or for a client.
type LabResult struct {
	id             uuid.UUID
	clientID       uuid.UUID
	nutritionistID uuid.UUID
	title          string
	resultType     string
	testDate       *time.Time
	filePath       string
	originalName   string
	fileType       string
	fileSize       int64
	link           *string
	notes          string
	createdAt      time.Time
}

// NewLabResult creates a new LabResult for persistence (ID assigned by DB).
func NewLabResult(clientID, nutritionistID uuid.UUID, title, resultType string, testDate *time.Time, notes string, link *string) *LabResult {
	return &LabResult{
		clientID:       clientID,
		nutritionistID: nutritionistID,
		title:          title,
		resultType:     resultType,
		testDate:       testDate,
		notes:          notes,
		link:           link,
	}
}

// ReconstituteLabResult rebuilds a LabResult from persisted storage data.
func ReconstituteLabResult(id, clientID, nutritionistID uuid.UUID, title, resultType string, testDate *time.Time, filePath, originalName, fileType string, fileSize int64, link *string, notes string, createdAt time.Time) *LabResult {
	return &LabResult{
		id:             id,
		clientID:       clientID,
		nutritionistID: nutritionistID,
		title:          title,
		resultType:     resultType,
		testDate:       testDate,
		filePath:       filePath,
		originalName:   originalName,
		fileType:       fileType,
		fileSize:       fileSize,
		link:           link,
		notes:          notes,
		createdAt:      createdAt,
	}
}

// Getters
func (r *LabResult) ID() uuid.UUID             { return r.id }
func (r *LabResult) ClientID() uuid.UUID       { return r.clientID }
func (r *LabResult) NutritionistID() uuid.UUID { return r.nutritionistID }
func (r *LabResult) Title() string             { return r.title }
func (r *LabResult) ResultType() string        { return r.resultType }
func (r *LabResult) TestDate() *time.Time      { return r.testDate }
func (r *LabResult) FilePath() string          { return r.filePath }
func (r *LabResult) OriginalName() string      { return r.originalName }
func (r *LabResult) FileType() string          { return r.fileType }
func (r *LabResult) FileSize() int64           { return r.fileSize }
func (r *LabResult) Link() *string             { return r.link }
func (r *LabResult) Notes() string             { return r.notes }
func (r *LabResult) CreatedAt() time.Time      { return r.createdAt }

// Setters for infrastructure and application layer use
func (r *LabResult) SetFilePath(v string)     { r.filePath = v }
func (r *LabResult) SetOriginalName(v string) { r.originalName = v }
func (r *LabResult) SetFileType(v string)     { r.fileType = v }
func (r *LabResult) SetFileSize(v int64)      { r.fileSize = v }
func (r *LabResult) SetLink(v *string)        { r.link = v }
