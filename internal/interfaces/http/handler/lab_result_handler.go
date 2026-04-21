package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appLabResult "github.com/ranjbar-dev/nutritrack/internal/application/labresult"
	labresultEntity "github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// LabResultHandler handles HTTP requests for lab result operations.
type LabResultHandler struct {
	service *appLabResult.LabResultService
}

// NewLabResultHandler creates a new LabResultHandler.
func NewLabResultHandler(service *appLabResult.LabResultService) *LabResultHandler {
	return &LabResultHandler{service: service}
}

// Upload handles POST /clients/:id/lab-results
// Accepts multipart/form-data with field "file" (PDF/JPEG/PNG, max 10 MB) and optional "notes".
func (h *LabResultHandler) Upload(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	// Parse multipart file — enforce 10 MB limit at the framework level
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1024*1024+512)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if fileHeader.Size > 10*1024*1024 {
		dto.Abort(c, shared.ErrFileTooLarge)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}
	defer file.Close()

	result, svcErr := h.service.UploadLabResult(
		c.Request.Context(),
		clientID,
		callerIDVal.(uuid.UUID),
		callerRoleVal.(string),
		file,
		fileHeader.Filename,
		fileHeader.Size,
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.Created(c, labResultToMap(result))
}

// List handles GET /clients/:id/lab-results
func (h *LabResultHandler) List(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	pg := dto.ParsePagination(c)

	results, total, svcErr := h.service.ListClientLabResults(
		c.Request.Context(),
		clientID,
		callerIDVal.(uuid.UUID),
		callerRoleVal.(string),
		pg.Limit(),
		pg.Offset(),
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	items := make([]map[string]any, len(results))
	for i, r := range results {
		items[i] = labResultToMap(r)
	}

	dto.Paginated(c, items, total, pg.Page, pg.PageSize)
}

// Download handles GET /lab-results/:id/download
func (h *LabResultHandler) Download(c *gin.Context) {
	labResultIDStr := c.Param("id")
	labResultID, err := uuid.Parse(labResultIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	result, svcErr := h.service.GetLabResultForDownload(
		c.Request.Context(),
		labResultID,
		callerIDVal.(uuid.UUID),
		callerRoleVal.(string),
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	// Serve the file as an attachment
	c.FileAttachment(result.FilePath, result.OriginalName)
}

func labResultToMap(r *labresultEntity.LabResult) map[string]any {
	return map[string]any{
		"id":              r.ID,
		"client_id":       r.ClientID,
		"nutritionist_id": r.NutritionistID,
		"original_name":   r.OriginalName,
		"file_type":       r.FileType,
		"file_size":       r.FileSize,
		"notes":           r.Notes,
		"created_at":      r.CreatedAt,
	}
}
