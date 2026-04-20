package handler

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

type TrackingHandler struct {
	svc        *service.TrackingService
	uploadsDir string
}

func NewTrackingHandler(svc *service.TrackingService, uploadsDir string) *TrackingHandler {
	return &TrackingHandler{svc: svc, uploadsDir: uploadsDir}
}

func (h *TrackingHandler) LogFood(c *gin.Context) {
	var req dto.LogFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.LogFood(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت وعده غذایی")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) ListFoodLogs(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListFoodLogs(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت وعده‌های غذایی")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) LogWater(c *gin.Context) {
	var req dto.LogWaterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.LogWater(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت مصرف آب")
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TrackingHandler) ListWaterLogs(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListWaterLogs(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت مصرف آب")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) UpsertSleep(c *gin.Context) {
	var req dto.UpsertSleepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.UpsertSleep(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت خواب")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) GetSleepLog(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.GetSleepLog(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت خواب")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) LogExercise(c *gin.Context) {
	var req dto.LogExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.LogExercise(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت تمرین")
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TrackingHandler) ListExerciseLogs(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListExerciseLogs(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت تمرین‌ها")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) LogMedication(c *gin.Context) {
	var req dto.LogMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.LogMedication(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت مصرف دارو")
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TrackingHandler) ListMedicationLogs(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListMedicationLogs(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت مصرف دارو")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) UpsertBodyMeasurement(c *gin.Context) {
	var req dto.UpsertBodyMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.UpsertBodyMeasurement(c.Request.Context(), clientID, clientID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت اندازه‌گیری")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) GetBodyMeasurement(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.GetBodyMeasurement(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت اندازه‌گیری")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) GetMeasurementHistory(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 90)
	resp, err := h.svc.ListBodyMeasurements(c.Request.Context(), clientID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت تاریخچه اندازه‌گیری")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) GetWeightHistory(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 90)
	resp, err := h.svc.GetWeightHistory(c.Request.Context(), clientID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت تاریخچه وزن")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) GetDailyDashboard(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.GetDailyDashboard(c.Request.Context(), clientID, service.NormalizeSingleDate(c.Query("date")))
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت خلاصه روزانه")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) UploadLabResult(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "فایل بیش از حد مجاز است"})
		return
	}
	var req dto.CreateLabResultRequest
	if err := c.ShouldBind(&req); err != nil {
		h.badRequest(c)
		return
	}
	var fileReader io.Reader
	var fileSize int64
	var originalFilename string
	if fh, err := c.FormFile("file"); err == nil {
		file, openErr := fh.Open()
		if openErr != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "فایل بارگذاری‌شده قابل خواندن نیست"})
			return
		}
		defer file.Close()
		fileReader = file
		fileSize = fh.Size
		originalFilename = fh.Filename
	}
	resp, err := h.svc.CreateLabResult(c.Request.Context(), clientID, req, fileReader, fileSize, originalFilename)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت نتیجه آزمایش")
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TrackingHandler) ListLabResults(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListLabResults(c.Request.Context(), clientID)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت نتایج آزمایش")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListFoodLogs(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 14)
	resp, err := h.svc.ListFoodLogsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت وعده‌های غذایی بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListWaterLogs(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 14)
	resp, err := h.svc.ListWaterLogsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت مصرف آب بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListSleepLogs(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 14)
	resp, err := h.svc.ListSleepLogsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت خواب بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListExerciseLogs(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 14)
	resp, err := h.svc.ListExerciseLogsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت تمرین‌های بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListMedicationLogs(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 14)
	resp, err := h.svc.ListMedicationLogsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت مصرف داروی بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListBodyMeasurements(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 90)
	resp, err := h.svc.ListBodyMeasurementsForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت اندازه‌گیری‌های بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriGetWeightHistory(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	from, to := service.NormalizeRange(c.Query("from"), c.Query("to"), 90)
	resp, err := h.svc.GetWeightHistoryForNutritionist(c.Request.Context(), clientID, nutritionistID, from, to)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت تاریخچه وزن بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriUpsertBodyMeasurement(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	var req dto.UpsertBodyMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.badRequest(c)
		return
	}
	resp, err := h.svc.UpsertBodyMeasurement(c.Request.Context(), clientID, nutritionistID, req)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در ثبت اندازه‌گیری بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriListLabResults(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	resp, err := h.svc.ListLabResultsForNutritionist(c.Request.Context(), clientID, nutritionistID)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت نتایج آزمایش بیمار")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TrackingHandler) NutriDownloadLabResult(c *gin.Context) {
	clientID, nutritionistID, ok := h.nutriClientAndOwner(c)
	if !ok {
		return
	}
	labID, err := uuid.Parse(c.Param("labId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return
	}
	resp, err := h.svc.GetLabResultForNutritionist(c.Request.Context(), labID, clientID, nutritionistID)
	if err != nil {
		h.handleTrackingError(c, err, "خطا در دریافت فایل آزمایش")
		return
	}
	if !resp.HasFile || resp.FilePath == nil || resp.OriginalFilename == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "این نتیجه آزمایش فایل ندارد"})
		return
	}
	absPath := filepath.Join(h.uploadsDir, *resp.FilePath)
	c.FileAttachment(absPath, *resp.OriginalFilename)
}

func (h *TrackingHandler) nutriClientAndOwner(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	nutritionistID, ok := authUUID(c)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	clientID, err := uuid.Parse(c.Param("clientId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
		return uuid.Nil, uuid.Nil, false
	}
	return clientID, nutritionistID, true
}

func authUUID(c *gin.Context) (uuid.UUID, bool) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return uuid.Nil, false
	}
	return userID, true
}

func (h *TrackingHandler) handleTrackingError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrTrackingNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrTrackingUnauthorized):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrLabFileMissing), errors.Is(err, service.ErrLabFileInvalidType), errors.Is(err, service.ErrLabFileTooLarge):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: fallback})
	}
}

func (h *TrackingHandler) badRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
}
