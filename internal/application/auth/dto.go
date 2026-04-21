package auth

import "github.com/google/uuid"

// LoginRequest for email/password login (superadmin + nutritionist).
type LoginRequest struct {
	Email    string
	Password string
}

// OTPSendRequest for client OTP flow.
type OTPSendRequest struct {
	Mobile string
}

// OTPVerifyRequest for verifying the OTP code.
type OTPVerifyRequest struct {
	Mobile string
	Code   string
}

// RefreshRequest for token refresh.
type RefreshRequest struct {
	RefreshToken string
}

// AuthResponse is returned after successful login/OTP verify/refresh.
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	UserID       uuid.UUID `json:"user_id"`
	Role         string    `json:"role"`
}
