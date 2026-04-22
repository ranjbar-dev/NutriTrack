package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrFoodRequestInvalidInput is returned when required fields are missing.
var ErrFoodRequestInvalidInput = errors.New("food request requires clientID, nutritionistID, and foodName")

type FoodRequestStatus string

const (
	FoodRequestStatusPending  FoodRequestStatus = "pending"
	FoodRequestStatusApproved FoodRequestStatus = "approved"
	FoodRequestStatusRejected FoodRequestStatus = "rejected"
)

// FoodRequest is the food-request aggregate. All fields are unexported;
// state is accessed through getters and mutated through behaviour methods.
type FoodRequest struct {
	id              uuid.UUID
	clientID        uuid.UUID
	nutritionistID  uuid.UUID
	foodName        string
	status          FoodRequestStatus
	rejectionReason *string
	createdFoodID   *uuid.UUID
	createdAt       time.Time
	updatedAt       time.Time
}

// NewFoodRequest creates a new, unreviewed FoodRequest aggregate.
// clientID, nutritionistID, and foodName are all required.
func NewFoodRequest(clientID, nutritionistID uuid.UUID, foodName string) (*FoodRequest, error) {
	if clientID == uuid.Nil || nutritionistID == uuid.Nil || foodName == "" {
		return nil, ErrFoodRequestInvalidInput
	}
	return &FoodRequest{
		clientID:       clientID,
		nutritionistID: nutritionistID,
		foodName:       foodName,
		status:         FoodRequestStatusPending,
	}, nil
}

// FromPersistence reconstructs a FoodRequest from stored data.
// This function is intended for use by the infrastructure layer only.
func FromPersistence(
	id, clientID, nutritionistID uuid.UUID,
	foodName string,
	status FoodRequestStatus,
	rejectionReason *string,
	createdFoodID *uuid.UUID,
	createdAt, updatedAt time.Time,
) *FoodRequest {
	return &FoodRequest{
		id:              id,
		clientID:        clientID,
		nutritionistID:  nutritionistID,
		foodName:        foodName,
		status:          status,
		rejectionReason: rejectionReason,
		createdFoodID:   createdFoodID,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// Hydrate populates DB-generated fields after a successful insert.
// Intended for use by the infrastructure layer only.
func (r *FoodRequest) Hydrate(id uuid.UUID, status FoodRequestStatus, createdAt, updatedAt time.Time) {
	r.id = id
	r.status = status
	r.createdAt = createdAt
	r.updatedAt = updatedAt
}

// Getters

func (r *FoodRequest) GetID() uuid.UUID              { return r.id }
func (r *FoodRequest) GetClientID() uuid.UUID        { return r.clientID }
func (r *FoodRequest) GetNutritionistID() uuid.UUID  { return r.nutritionistID }
func (r *FoodRequest) GetFoodName() string           { return r.foodName }
func (r *FoodRequest) GetStatus() FoodRequestStatus  { return r.status }
func (r *FoodRequest) GetRejectionReason() *string   { return r.rejectionReason }
func (r *FoodRequest) GetCreatedFoodID() *uuid.UUID  { return r.createdFoodID }
func (r *FoodRequest) GetCreatedAt() time.Time       { return r.createdAt }
func (r *FoodRequest) GetUpdatedAt() time.Time       { return r.updatedAt }

// State predicates

func (r *FoodRequest) IsPending() bool  { return r.status == FoodRequestStatusPending }
func (r *FoodRequest) IsApproved() bool { return r.status == FoodRequestStatusApproved }
func (r *FoodRequest) IsRejected() bool { return r.status == FoodRequestStatusRejected }
