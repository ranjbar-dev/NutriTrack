package entity

import (
	"time"

	"github.com/google/uuid"
)

// Medication is the medication/supplement domain aggregate.
type Medication struct {
	id             uuid.UUID
	name           string
	nameNormalized string
	description    string
	unit           string
	createdBy      *uuid.UUID
	isActive       bool
	createdAt      time.Time
	updatedAt      time.Time
}

// NewMedication creates a new Medication for persistence (ID assigned by DB).
func NewMedication(name, nameNormalized, description, unit string, createdBy *uuid.UUID) *Medication {
	return &Medication{
		name:           name,
		nameNormalized: nameNormalized,
		description:    description,
		unit:           unit,
		createdBy:      createdBy,
		isActive:       true,
	}
}

// ReconstituteMedication rebuilds a Medication from persisted storage data.
func ReconstituteMedication(id uuid.UUID, name, nameNormalized, description, unit string, createdBy *uuid.UUID, isActive bool, createdAt, updatedAt time.Time) *Medication {
	return &Medication{
		id:             id,
		name:           name,
		nameNormalized: nameNormalized,
		description:    description,
		unit:           unit,
		createdBy:      createdBy,
		isActive:       isActive,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// Getters
func (m *Medication) ID() uuid.UUID          { return m.id }
func (m *Medication) Name() string           { return m.name }
func (m *Medication) NameNormalized() string { return m.nameNormalized }
func (m *Medication) Description() string    { return m.description }
func (m *Medication) Unit() string           { return m.unit }
func (m *Medication) CreatedBy() *uuid.UUID  { return m.createdBy }
func (m *Medication) IsActive() bool         { return m.isActive }
func (m *Medication) CreatedAt() time.Time   { return m.createdAt }
func (m *Medication) UpdatedAt() time.Time   { return m.updatedAt }

// Setters for infrastructure and application layer use
func (m *Medication) SetID(id uuid.UUID)         { m.id = id }
func (m *Medication) SetName(name string)        { m.name = name }
func (m *Medication) SetNameNormalized(n string) { m.nameNormalized = n }
func (m *Medication) SetDescription(d string)    { m.description = d }
func (m *Medication) SetUnit(u string)           { m.unit = u }
func (m *Medication) SetIsActive(v bool)         { m.isActive = v }
func (m *Medication) SetCreatedAt(t time.Time)   { m.createdAt = t }
func (m *Medication) SetUpdatedAt(t time.Time)   { m.updatedAt = t }
