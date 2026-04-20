package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// PushHandler handles push notification subscription and preference endpoints.
type PushHandler struct {
	notifSvc service.NotificationService
}

// NewPushHandler creates a new PushHandler.
func NewPushHandler(notifSvc service.NotificationService) *PushHandler {
	return &PushHandler{notifSvc: notifSvc}
}

// Subscribe registers a push subscription for the authenticated client.
func (h *PushHandler) Subscribe(c *gin.Context) {
	clientID := c.GetString("userID")
	var req dto.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserAgent = c.Request.UserAgent()
	if err := h.notifSvc.RegisterSubscription(c.Request.Context(), clientID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ثبت اشتراک ناموفق بود"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "اشتراک ثبت شد"})
}

// Unsubscribe removes a push subscription for the authenticated client.
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	clientID := c.GetString("userID")
	var req dto.UnsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.notifSvc.RemoveSubscription(c.Request.Context(), clientID, req.Endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "لغو اشتراک ناموفق بود"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "اشتراک حذف شد"})
}

// GetPreferences returns notification preferences for the authenticated client.
func (h *PushHandler) GetPreferences(c *gin.Context) {
	clientID := c.GetString("userID")
	prefs, err := h.notifSvc.GetPreferences(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در دریافت تنظیمات"})
		return
	}
	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences saves notification preferences for the authenticated client.
func (h *PushHandler) UpdatePreferences(c *gin.Context) {
	clientID := c.GetString("userID")
	var req dto.NotificationPreferencesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.notifSvc.UpdatePreferences(c.Request.Context(), clientID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ذخیره تنظیمات"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
