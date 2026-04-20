package dto

import "time"

// ─── Message DTOs ───────────────────────────────────────────────────────────

type SendMessageRequest struct {
	Content        *string `json:"content" binding:"omitempty,max=5000"`
	AttachmentType *string `json:"attachment_type" binding:"omitempty,oneof=image pdf"`
}

type MessageResponse struct {
	ID             string  `json:"id"`
	SenderID       string  `json:"sender_id"`
	ReceiverID     string  `json:"receiver_id"`
	Content        *string `json:"content"`
	AttachmentType *string `json:"attachment_type"`
	AttachmentPath *string `json:"attachment_path,omitempty"`
	AttachmentName *string `json:"attachment_name"`
	SentAt         time.Time `json:"sent_at"`
	ReadAt         *time.Time `json:"read_at"`
}

type UnreadCountResponse struct {
	Count int32 `json:"count"`
}

// ─── Food Request DTOs ───────────────────────────────────────────────────────

type FoodRequestCreateRequest struct {
	FoodName    string  `json:"food_name" binding:"required,max=200"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}

type FoodRequestRejectRequest struct {
	RejectionReason *string `json:"rejection_reason" binding:"omitempty,max=1000"`
}

type FoodRequestResponse struct {
	ID              string     `json:"id"`
	FoodName        string     `json:"food_name"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	RejectionReason *string    `json:"rejection_reason"`
	RequestedBy     string     `json:"requested_by"`
	ReviewedBy      *string    `json:"reviewed_by"`
	ClientName      *string    `json:"client_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ─── Client Management DTOs ──────────────────────────────────────────────────

type ClientListItem struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Mobile    *string   `json:"mobile"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ClientListResponse struct {
	Items []ClientListItem `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
}

type ClientProfileResponse struct {
	ID             string     `json:"id"`
	FullName       string     `json:"full_name"`
	Email          *string    `json:"email"`
	Mobile         *string    `json:"mobile"`
	DateOfBirth    *string    `json:"date_of_birth"`
	HeightCm       *float32   `json:"height_cm"`
	Gender         *string    `json:"gender"`
	NutritionistID *string    `json:"nutritionist_id"`
	IsActive       bool       `json:"is_active"`
	Notes          *string    `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UpdateClientProfileRequest struct {
	DateOfBirth *string  `json:"date_of_birth" binding:"omitempty"`
	HeightCm    *float32 `json:"height_cm" binding:"omitempty,min=50,max=250"`
}
