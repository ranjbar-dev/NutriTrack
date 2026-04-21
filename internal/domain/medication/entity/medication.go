package entity

import (
	"time"

	"github.com/google/uuid"
)

// Medication is the medication/supplement domain entity.
type Medication struct {
	ID             uuid.UUID
	Name           string
	NameNormalized string
	Description    string
	Unit           string
	CreatedBy      *uuid.UUID
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
