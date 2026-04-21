package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/application/auth"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/middleware"
)

// AuthHandler handles all authentication-related HTTP endpoints.
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given AuthService.
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /api/v1/auth/login (superadmin + nutritionist).
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), auth.LoginRequest{
		Email:    strings.TrimSpace(strings.ToLower(req.Email)),
		Password: req.Password,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, resp)
}

// SendOTP handles POST /api/v1/auth/otp/send (client).
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if err := h.authService.SendOTP(c.Request.Context(), auth.OTPSendRequest{
		Mobile: strings.TrimSpace(req.Mobile),
	}); err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, gin.H{"message": "کد تأیید ارسال شد"})
}

// VerifyOTP handles POST /api/v1/auth/otp/verify (client).
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
		Code   string `json:"code"   binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	resp, err := h.authService.VerifyOTP(c.Request.Context(), auth.OTPVerifyRequest{
		Mobile: strings.TrimSpace(req.Mobile),
		Code:   strings.TrimSpace(req.Code),
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, resp)
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), auth.RefreshRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		appErr, ok := err.(*shared.AppError)
		if !ok {
			appErr = shared.ErrInternal
		}
		dto.Abort(c, appErr)
		return
	}

	dto.OK(c, resp)
}

// Logout handles POST /api/v1/auth/logout (authenticated).
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Abort(c, shared.ErrValidation)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		dto.Abort(c, shared.ErrInternal)
		return
	}

	// Access token expiry is short (15 min), so we don't blacklist it here.
	_ = c.GetHeader("Authorization")

	c.JSON(http.StatusOK, gin.H{"message": "با موفقیت خارج شدید"})
}

// Me handles GET /api/v1/auth/me (authenticated — returns current user info).
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get(middleware.AuthUserIDKey)
	userRole, _ := c.Get(middleware.AuthUserRoleKey)

	dto.OK(c, gin.H{
		"user_id": userID,
		"role":    userRole,
	})
}
