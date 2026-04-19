package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given auth service.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /api/auth/login (AUTH-01, AUTH-03).
// Authenticates admin and nutritionist users with email/password.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	tokens, userResp, err := h.authService.LoginWithPassword(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)
	c.JSON(http.StatusOK, dto.AuthResponse{User: *userResp})
}

// RequestOTP handles POST /api/auth/otp/request (AUTH-05).
// Always returns 200 regardless of user existence to prevent enumeration (T-04-02).
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req dto.OTPRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شماره موبایل نامعتبر است"})
		return
	}

	// Ignore error — always return success for anti-enumeration
	_ = h.authService.RequestOTP(c.Request.Context(), req.Mobile)
	c.JSON(http.StatusOK, gin.H{"message": "کد تایید ارسال شد"})
}

// VerifyOTP handles POST /api/auth/otp/verify (AUTH-05, AUTH-06).
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.OTPVerifyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
		return
	}

	tokens, userResp, err := h.authService.VerifyOTP(c.Request.Context(), req.Mobile, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)
	c.JSON(http.StatusOK, dto.AuthResponse{User: *userResp})
}

// Refresh handles POST /api/auth/refresh (AUTH-08).
// Rotates refresh token and issues a new token pair.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	tokens, err := h.authService.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "توکن با موفقیت تمدید شد"})
}

// Logout handles POST /api/auth/logout.
// Revokes the refresh token and clears auth cookies.
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	_ = h.authService.Logout(c.Request.Context(), refreshToken)

	clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "خروج با موفقیت انجام شد"})
}

// GetMe handles GET /api/auth/me (session persistence on page refresh).
// Returns the current user's profile from a valid access_token cookie.
// Frontend checkAuth() calls this on hard refresh to restore session state.
func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
		return
	}

	userResp, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "کاربر یافت نشد"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{User: *userResp})
}

// setAuthCookies sets httpOnly secure cookies for access and refresh tokens (D-01).
// access_token: path=/api, maxAge=900 (15min), secure=true, httpOnly=true
// refresh_token: path=/api/auth/refresh, maxAge=2592000 (30d), secure=true, httpOnly=true
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie("access_token", accessToken, 900, "/api", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 2592000, "/api/auth/refresh", "", true, true)
}

// clearAuthCookies removes auth cookies by setting maxAge to -1.
func clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/api", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/api/auth/refresh", "", true, true)
}
