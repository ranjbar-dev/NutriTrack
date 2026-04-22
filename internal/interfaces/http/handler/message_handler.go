package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appMessage "github.com/ranjbar-dev/nutritrack/internal/application/message"
	appPush "github.com/ranjbar-dev/nutritrack/internal/application/push"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// MessageHandler handles HTTP requests for chat messages.
type MessageHandler struct {
	service *appMessage.MessageService
	pushSvc *appPush.PushService
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(service *appMessage.MessageService, pushSvc *appPush.PushService) *MessageHandler {
	return &MessageHandler{service: service, pushSvc: pushSvc}
}

// SendAsClient handles POST /messages — client sends to nutritionist.
// Accepts multipart/form-data with "content" (text) and optional "file" attachment.
func (h *MessageHandler) SendAsClient(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "client" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	content := c.PostForm("content")

	var (
		attachReader io.Reader
		attachName   string
		attachSize   int64
	)

	fileHeader, formErr := c.FormFile("file")
	if formErr == nil {
		// Enforce max size at HTTP layer
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1024*1024+512)
		if fileHeader.Size > 10*1024*1024 {
			dto.Abort(c, shared.ErrFileTooLarge)
			return
		}
		f, openErr := fileHeader.Open()
		if openErr != nil {
			dto.Abort(c, shared.ErrInternal)
			return
		}
		defer f.Close()
		attachReader = f
		attachName = fileHeader.Filename
		attachSize = fileHeader.Size
	}

	result, svcErr := h.service.SendAsClient(
		c.Request.Context(),
		callerIDVal.(uuid.UUID),
		content,
		attachReader,
		attachName,
		attachSize,
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	if h.pushSvc != nil {
		receiverID := result.ReceiverID()
		go func() {
			_ = h.pushSvc.Send(context.Background(), receiverID, "پیام جدید", "یک پیام جدید دریافت کردید")
		}()
	}

	dto.Created(c, appMessage.MapMessageResponse(result, callerIDVal.(uuid.UUID)))
}

// SendAsNutritionist handles POST /clients/:id/messages — nutritionist sends to client.
func (h *MessageHandler) SendAsNutritionist(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "nutritionist" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	content := c.PostForm("content")

	var (
		attachReader io.Reader
		attachName   string
		attachSize   int64
	)

	fileHeader, formErr := c.FormFile("file")
	if formErr == nil {
		if fileHeader.Size > 10*1024*1024 {
			dto.Abort(c, shared.ErrFileTooLarge)
			return
		}
		f, openErr := fileHeader.Open()
		if openErr != nil {
			dto.Abort(c, shared.ErrInternal)
			return
		}
		defer f.Close()
		attachReader = f
		attachName = fileHeader.Filename
		attachSize = fileHeader.Size
	}

	result, svcErr := h.service.SendAsNutritionist(
		c.Request.Context(),
		callerIDVal.(uuid.UUID),
		clientID,
		content,
		attachReader,
		attachName,
		attachSize,
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	if h.pushSvc != nil {
		receiverID := result.ReceiverID()
		go func() {
			_ = h.pushSvc.Send(context.Background(), receiverID, "پیام جدید", "متخصص تغذیه شما پیام جدید فرستاد")
		}()
	}

	dto.Created(c, appMessage.MapMessageResponse(result, callerIDVal.(uuid.UUID)))
}

// GetClientMessages handles GET /messages — client views their conversation with nutritionist.
func (h *MessageHandler) GetClientMessages(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "client" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	pg := dto.ParsePagination(c)
	msgs, total, svcErr := h.service.GetClientConversation(
		c.Request.Context(),
		callerIDVal.(uuid.UUID),
		int32(pg.Limit()),
		int32(pg.Offset()),
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	items := make([]map[string]any, len(msgs))
	for i, m := range msgs {
		items[i] = appMessage.MapMessageResponse(m, callerIDVal.(uuid.UUID))
	}
	dto.Paginated(c, items, total, pg.Page, pg.PageSize)
}

// GetNutritionistMessages handles GET /clients/:id/messages — nutritionist views conversation with client.
func (h *MessageHandler) GetNutritionistMessages(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)

	if callerRoleVal.(string) != "nutritionist" {
		dto.Abort(c, shared.ErrForbidden)
		return
	}

	pg := dto.ParsePagination(c)
	msgs, total, svcErr := h.service.GetNutritionistConversation(
		c.Request.Context(),
		callerIDVal.(uuid.UUID),
		clientID,
		int32(pg.Limit()),
		int32(pg.Offset()),
	)
	if svcErr != nil {
		if appErr, ok := svcErr.(*shared.AppError); ok {
			dto.Abort(c, appErr)
			return
		}
		dto.Abort(c, shared.ErrInternal)
		return
	}

	items := make([]map[string]any, len(msgs))
	for i, m := range msgs {
		items[i] = appMessage.MapMessageResponse(m, callerIDVal.(uuid.UUID))
	}
	dto.Paginated(c, items, total, pg.Page, pg.PageSize)
}

// GetUnreadCount handles GET /messages/unread-count — any authenticated user.
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
	count, svcErr := h.service.GetUnreadCount(c.Request.Context(), callerIDVal.(uuid.UUID))
	if svcErr != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}
	dto.OK(c, map[string]any{"unread_count": count})
}
