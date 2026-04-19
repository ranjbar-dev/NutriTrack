package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// MedicationHandler handles medication CRUD HTTP endpoints.
type MedicationHandler struct {
	medService *service.MedicationService
}

// NewMedicationHandler creates a new MedicationHandler with the given medication service.
func NewMedicationHandler(medService *service.MedicationService) *MedicationHandler {
	return &MedicationHandler{medService: medService}
}

// Create handles POST /api/medications.
func (h *MedicationHandler) Create(c *gin.Context) {
	var req dto.CreateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.medService.CreateMedication(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMedicationInvalidName):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrMedicationDuplicate):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ثبت دارو"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Get handles GET /api/medications/:id.
func (h *MedicationHandler) Get(c *gin.Context) {
	medID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	resp, err := h.medService.GetMedication(c.Request.Context(), medID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMedicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت دارو"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List handles GET /api/medications.
func (h *MedicationHandler) List(c *gin.Context) {
	var query dto.MedicationListQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	resp, err := h.medService.ListMedications(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت لیست داروها"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update handles PUT /api/medications/:id.
func (h *MedicationHandler) Update(c *gin.Context) {
	medID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	var req dto.UpdateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	resp, err := h.medService.UpdateMedication(c.Request.Context(), medID, userID, c.GetString("role"), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMedicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrMedicationInvalidName):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrMedicationUnauthorizedEdit):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrMedicationDuplicate):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ویرایش دارو"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Delete handles DELETE /api/medications/:id.
func (h *MedicationHandler) Delete(c *gin.Context) {
	medID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	err = h.medService.DeleteMedication(c.Request.Context(), medID, userID, c.GetString("role"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMedicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrMedicationUnauthorizedDelete):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در حذف دارو"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "دارو با موفقیت حذف شد"})
}
