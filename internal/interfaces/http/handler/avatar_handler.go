package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appUser "github.com/ranjbar-dev/nutritrack/internal/application/user"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

const maxAvatarSize = 5 * 1024 * 1024 // 5 MB

// AvatarHandler handles profile picture upload.
type AvatarHandler struct {
	avatarSvc *appUser.AvatarService
}

// NewAvatarHandler constructs an AvatarHandler.
func NewAvatarHandler(svc *appUser.AvatarService) *AvatarHandler {
	return &AvatarHandler{avatarSvc: svc}
}

// Upload handles PUT /users/:id/avatar
func (h *AvatarHandler) Upload(c *gin.Context) {
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	callerIDRaw, _ := c.Get(middleware.AuthUserIDKey)
	callerRoleRaw, _ := c.Get(middleware.AuthUserRoleKey)
	callerID, _ := callerIDRaw.(uuid.UUID)
	callerRole, _ := callerRoleRaw.(string)

	// Limit request body size to prevent memory exhaustion
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAvatarSize+512)

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		dto.Abort(c, shared.ErrValidation.WithMessage("فایل آواتار ارسال نشده است"))
		return
	}

	if fileHeader.Size > maxAvatarSize {
		dto.Abort(c, shared.ErrFileTooLarge)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}
	defer f.Close()

	// Read magic bytes (first 12 bytes) for MIME validation
	magicBuf := make([]byte, 12)
	n, _ := io.ReadFull(f, magicBuf)
	magicBuf = magicBuf[:n]

	// Reconstruct full reader: prepend magic bytes back so the storage writer
	// receives the complete file from byte 0
	fullReader := io.MultiReader(bytes.NewReader(magicBuf), f)

	user, err := h.avatarSvc.UploadAvatar(
		c.Request.Context(),
		targetID,
		callerID,
		callerRole,
		magicBuf,
		fullReader,
		maxAvatarSize,
	)
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, gin.H{
		"avatar_url": user.GetAvatarURL(),
		"message":    "تصویر پروفایل با موفقیت بارگذاری شد",
	})
}
