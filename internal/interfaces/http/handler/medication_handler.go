package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appMed "github.com/ranjbar-dev/nutritrack/internal/application/medication"
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// MedicationHandler handles medication CRUD and search endpoints.
type MedicationHandler struct {
	svc *appMed.MedicationService
}

// NewMedicationHandler constructs a MedicationHandler.
func NewMedicationHandler(svc *appMed.MedicationService) *MedicationHandler {
	return &MedicationHandler{svc: svc}
}

// medicationRequest is the shared request body for POST and PATCH.
type medicationRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
	Unit        string `json:"unit"        binding:"required"`
}

// Create handles POST /api/v1/medications.
func (h *MedicationHandler) Create(c *gin.Context) {
	var req medicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	callerID, _ := callerIDRaw.(uuid.UUID)
	callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
	callerRole, _ := callerRoleRaw.(string)

	med, err := h.svc.CreateMedication(c.Request.Context(), appMed.CreateMedicationRequest{
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		CallerID:    callerID,
		CallerRole:  callerRole,
	})
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.Created(c, toMedicationResponse(med))
}

// GetOne handles GET /api/v1/medications/:id.
func (h *MedicationHandler) GetOne(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	med, err := h.svc.GetMedication(c.Request.Context(), id)
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.OK(c, toMedicationResponse(med))
}

// Search handles GET /api/v1/medications.
func (h *MedicationHandler) Search(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	pg := dto.ParsePagination(c)

	meds, total, err := h.svc.SearchMedications(c.Request.Context(), query, pg.Limit(), pg.Offset())
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	resp := make([]gin.H, len(meds))
	for i, m := range meds {
		resp[i] = toMedicationResponse(m)
	}
	dto.Paginated(c, resp, total, pg.Page, pg.PageSize)
}

// Update handles PATCH /api/v1/medications/:id.
func (h *MedicationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	var req medicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	callerID, _ := callerIDRaw.(uuid.UUID)
	callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
	callerRole, _ := callerRoleRaw.(string)

	med, err := h.svc.UpdateMedication(c.Request.Context(), appMed.UpdateMedicationRequest{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
		CallerID:    callerID,
		CallerRole:  callerRole,
	})
	if err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	dto.OK(c, toMedicationResponse(med))
}

// Delete handles DELETE /api/v1/medications/:id.
func (h *MedicationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	callerID, _ := callerIDRaw.(uuid.UUID)
	callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
	callerRole, _ := callerRoleRaw.(string)

	if err := h.svc.DeleteMedication(c.Request.Context(), id, callerID, callerRole); err != nil {
		dto.Abort(c, toAppError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// --- helpers ---

// toMedicationResponse converts a domain Medication to a JSON-serialisable map.
func toMedicationResponse(m *entity.Medication) gin.H {
	resp := gin.H{
		"id":          m.ID,
		"name":        m.Name,
		"description": m.Description,
		"unit":        m.Unit,
		"is_active":   m.IsActive,
		"created_at":  m.CreatedAt,
		"updated_at":  m.UpdatedAt,
	}
	if m.CreatedBy != nil {
		resp["created_by"] = *m.CreatedBy
	} else {
		resp["created_by"] = nil
	}
	return resp
}
