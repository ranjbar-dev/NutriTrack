package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appNotif "github.com/ranjbar-dev/nutritrack/internal/application/notification"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// NotificationHandler handles notification preference endpoints.
type NotificationHandler struct {
	service *appNotif.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(svc *appNotif.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: svc}
}

// UpdatePreferences handles PATCH /notifications/preferences — upsert notification preferences.
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	callerID := c.MustGet(middleware.AuthUserIDKey).(uuid.UUID)

	var body struct {
		MealReminders  bool `json:"meal_reminders"`
		WaterReminders bool `json:"water_reminders"`
		MessageAlerts  bool `json:"message_alerts"`
		DietUpdates    bool `json:"diet_updates"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	pref, err := h.service.UpdatePreferences(c.Request.Context(), callerID, appNotif.UpdatePreferencesRequest{
		MealReminders:  body.MealReminders,
		WaterReminders: body.WaterReminders,
		MessageAlerts:  body.MessageAlerts,
		DietUpdates:    body.DietUpdates,
	})
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.OK(c, map[string]any{
		"id":              pref.ID,
		"user_id":         pref.UserID,
		"meal_reminders":  pref.MealReminders,
		"water_reminders": pref.WaterReminders,
		"message_alerts":  pref.MessageAlerts,
		"diet_updates":    pref.DietUpdates,
	})
}

// GetPreferences handles GET /notifications/preferences — retrieve notification preferences.
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	callerID := c.MustGet(middleware.AuthUserIDKey).(uuid.UUID)

	pref, err := h.service.GetPreferences(c.Request.Context(), callerID)
	if err != nil {
		dto.Abort(c, shared.ErrNotificationPreferenceNotFound)
		return
	}

	dto.OK(c, map[string]any{
		"id":              pref.ID,
		"user_id":         pref.UserID,
		"meal_reminders":  pref.MealReminders,
		"water_reminders": pref.WaterReminders,
		"message_alerts":  pref.MessageAlerts,
		"diet_updates":    pref.DietUpdates,
	})
}
