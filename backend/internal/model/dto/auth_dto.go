package dto

// LoginRequest represents the email/password login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// OTPRequestDTO represents an OTP request for client login.
type OTPRequestDTO struct {
	Mobile string `json:"mobile" binding:"required,iranian_mobile"`
}

// OTPVerifyDTO represents an OTP verification request.
type OTPVerifyDTO struct {
	Mobile string `json:"mobile" binding:"required,iranian_mobile"`
	Code   string `json:"code" binding:"required,len=6,numeric"`
}

// CreateNutritionistRequest represents the request to create a new nutritionist account.
type CreateNutritionistRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterClientRequest represents the request to register a new client.
type RegisterClientRequest struct {
	FullName    string   `json:"full_name" binding:"required,min=2,max=255"`
	Mobile      string   `json:"mobile" binding:"required,iranian_mobile"`
	DateOfBirth *string  `json:"date_of_birth" binding:"omitempty"`
	HeightCM    *float32 `json:"height_cm" binding:"omitempty,gt=0,lt=300"`
	Gender      *string  `json:"gender" binding:"omitempty,oneof=male female"`
	Notes       *string  `json:"notes" binding:"omitempty,max=1000"`
}

// AuthResponse represents the response after successful authentication.
type AuthResponse struct {
	User UserResponse `json:"user"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
}

// ErrorResponse represents an error returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents the health check response body.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
