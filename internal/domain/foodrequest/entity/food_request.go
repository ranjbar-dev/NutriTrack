package entity

import (
	"time"

	"github.com/google/uuid"
)

type FoodRequestStatus string

const (
	FoodRequestStatusPending  FoodRequestStatus = "pending"
	FoodRequestStatusApproved FoodRequestStatus = "approved"
	FoodRequestStatusRejected FoodRequestStatus = "rejected"
)

type FoodRequest struct {
	ID              uuid.UUID
	ClientID        uuid.UUID
	NutritionistID  uuid.UUID
	FoodName        string
	Status          FoodRequestStatus
	RejectionReason *string
	CreatedFoodID   *uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (r *FoodRequest) IsPending() bool  { return r.Status == FoodRequestStatusPending }
func (r *FoodRequest) IsApproved() bool { return r.Status == FoodRequestStatusApproved }
func (r *FoodRequest) IsRejected() bool { return r.Status == FoodRequestStatusRejected }
