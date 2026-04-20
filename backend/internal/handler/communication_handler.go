package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// CommunicationHandler handles messaging and food-request HTTP endpoints.
type CommunicationHandler struct {
	commService *service.CommunicationService
	uploadsDir  string
}

// NewCommunicationHandler creates a new CommunicationHandler.
func NewCommunicationHandler(commService *service.CommunicationService, uploadsDir string) *CommunicationHandler {
	return &CommunicationHandler{commService: commService, uploadsDir: uploadsDir}
}

// ─── Messaging ───────────────────────────────────────────────────────────────

// SendMessage handles POST /api/messages
// Accepts multipart/form-data with optional file attachment plus optional JSON content field.
func (h *CommunicationHandler) SendMessage(c *gin.Context) {
	senderID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	receiverIDStr := c.PostForm("receiver_id")
	receiverID, err := uuid.Parse(receiverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه گیرنده نامعتبر است"})
		return
	}

	content := c.PostForm("content")
	var contentPtr *string
	if content != "" {
		contentPtr = &content
	}

	var (
		fileReader interface{ Read([]byte) (int, error) } = nil
		fileSize   int64
		filename   string
		mimeStr    string
	)

	file, fh, ferr := c.Request.FormFile("attachment")
	if ferr == nil {
		defer file.Close()
		fileReader = file
		fileSize = fh.Size
		filename = fh.Filename
		mimeStr = fh.Header.Get("Content-Type")
	}

	var fr interface {
		Read([]byte) (int, error)
	} = fileReader

	msg, svcErr := h.commService.SendMessageTo(c.Request.Context(), senderID, receiverID, contentPtr, fr, fileSize, filename, mimeStr)
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.JSON(http.StatusCreated, msg)
}

// ListMessages handles GET /api/messages/:partnerId
func (h *CommunicationHandler) ListMessages(c *gin.Context) {
	userID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	partnerID, err := uuid.Parse(c.Param("partnerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه طرف مکالمه نامعتبر است"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

	msgs, svcErr := h.commService.GetMessages(c.Request.Context(), userID, partnerID, int32(limit), int32(offset))
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, msgs)
}

// PollNewMessages handles GET /api/messages/:partnerId/poll?since=<RFC3339>
func (h *CommunicationHandler) PollNewMessages(c *gin.Context) {
	userID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	partnerID, err := uuid.Parse(c.Param("partnerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه طرف مکالمه نامعتبر است"})
		return
	}

	sinceStr := c.Query("since")
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "فرمت زمان نامعتبر است — RFC3339 مورد نیاز است"})
		return
	}

	msgs, svcErr := h.commService.GetNewMessages(c.Request.Context(), userID, partnerID, since)
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, msgs)
}

// MarkRead handles PATCH /api/messages/:partnerId/read
func (h *CommunicationHandler) MarkRead(c *gin.Context) {
	receiverID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	senderID, err := uuid.Parse(c.Param("partnerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه فرستنده نامعتبر است"})
		return
	}

	if err := h.commService.MarkRead(c.Request.Context(), receiverID, senderID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در علامت‌گذاری پیام‌ها"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "پیام‌ها خوانده شدند"})
}

// GetUnreadCount handles GET /api/messages/unread-count
// IMPORTANT: must be registered BEFORE /:partnerId wildcard.
func (h *CommunicationHandler) GetUnreadCount(c *gin.Context) {
	userID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	resp, err := h.commService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت تعداد پیام‌های خوانده نشده"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DownloadAttachment handles GET /api/messages/attachment/:messageId
// IMPORTANT: must be registered BEFORE /:partnerId wildcard.
func (h *CommunicationHandler) DownloadAttachment(c *gin.Context) {
	userID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	messageID, err := uuid.Parse(c.Param("messageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه پیام نامعتبر است"})
		return
	}

	absPath, origName, svcErr := h.commService.GetMessageAttachment(c.Request.Context(), userID, messageID)
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.FileAttachment(absPath, origName)
}

// ─── Food Requests ────────────────────────────────────────────────────────────

// ClientCreateFoodRequest handles POST /api/client/food-requests
func (h *CommunicationHandler) ClientCreateFoodRequest(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	var req dto.FoodRequestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	fr, err := h.commService.CreateFoodRequest(c.Request.Context(), clientID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ثبت درخواست"})
		return
	}
	c.JSON(http.StatusCreated, fr)
}

// ClientListFoodRequests handles GET /api/client/food-requests
func (h *CommunicationHandler) ClientListFoodRequests(c *gin.Context) {
	clientID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	frs, err := h.commService.ListClientFoodRequests(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت درخواست‌ها"})
		return
	}
	c.JSON(http.StatusOK, frs)
}

// NutriListFoodRequests handles GET /api/nutritionist/food-requests
func (h *CommunicationHandler) NutriListFoodRequests(c *gin.Context) {
	nutriID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	frs, err := h.commService.ListNutriPendingFoodRequests(c.Request.Context(), nutriID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در دریافت درخواست‌ها"})
		return
	}
	c.JSON(http.StatusOK, frs)
}

// NutriApproveFoodRequest handles PATCH /api/nutritionist/food-requests/:requestId/approve
func (h *CommunicationHandler) NutriApproveFoodRequest(c *gin.Context) {
	nutriID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	requestID, err := uuid.Parse(c.Param("requestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه درخواست نامعتبر است"})
		return
	}

	fr, svcErr := h.commService.ApproveFoodRequest(c.Request.Context(), requestID, nutriID)
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, fr)
}

// NutriRejectFoodRequest handles PATCH /api/nutritionist/food-requests/:requestId/reject
func (h *CommunicationHandler) NutriRejectFoodRequest(c *gin.Context) {
	nutriID, ok := authUUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	requestID, err := uuid.Parse(c.Param("requestId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه درخواست نامعتبر است"})
		return
	}

	var req dto.FoodRequestRejectRequest
	_ = c.ShouldBindJSON(&req) // optional body

	fr, svcErr := h.commService.RejectFoodRequest(c.Request.Context(), requestID, nutriID, req.RejectionReason)
	if svcErr != nil {
		h.handleCommError(c, svcErr)
		return
	}
	c.JSON(http.StatusOK, fr)
}

// ─── error helper ─────────────────────────────────────────────────────────────

func (h *CommunicationHandler) handleCommError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrCommNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrCommUnauthorized):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrMsgAttachmentTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrMsgAttachmentInvalidType):
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrMsgNoContent):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrFoodRequestAlreadyReviewed):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطای داخلی سرور"})
	}
}
