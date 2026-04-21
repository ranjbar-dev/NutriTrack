package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appPush "github.com/ranjbar-dev/nutritrack/internal/application/push"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// PushHandler handles Web Push subscription management.
type PushHandler struct {
	service *appPush.PushService
}

// NewPushHandler creates a new PushHandler.
func NewPushHandler(service *appPush.PushService) *PushHandler {
	return &PushHandler{service: service}
}

// Subscribe handles POST /push/subscribe — register or update a push subscription.
func (h *PushHandler) Subscribe(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)

	var body struct {
		Endpoint string `json:"endpoint"`
		P256dh   string `json:"p256dh"`
		Auth     string `json:"auth"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Endpoint == "" || body.P256dh == "" || body.Auth == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	sub, err := h.service.Subscribe(c.Request.Context(), callerIDVal.(uuid.UUID), body.Endpoint, body.P256dh, body.Auth)
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.Created(c, map[string]any{
		"id":         sub.ID,
		"user_id":    sub.UserID,
		"endpoint":   sub.Endpoint,
		"created_at": sub.CreatedAt,
	})
}

// Unsubscribe handles DELETE /push/subscribe — remove a push subscription.
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)

	var body struct {
		Endpoint string `json:"endpoint"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Endpoint == "" {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if err := h.service.Unsubscribe(c.Request.Context(), callerIDVal.(uuid.UUID), body.Endpoint); err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	dto.NoContent(c)
}
