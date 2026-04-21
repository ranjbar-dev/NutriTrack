package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role constants
const (
	RoleSuperAdmin   = "superadmin"
	RoleNutritionist = "nutritionist"
	RoleClient       = "client"
)

// Gender constants
const (
	GenderMale   = "male"
	GenderFemale = "female"
)

// User is the root aggregate for identity and profile data.
type User struct {
	ID             uuid.UUID
	Role           string
	Mobile         string     // Iranian mobile, E.164 format: +989xxxxxxxxx
	Email          string     // nullable — only superadmin/nutritionist
	PasswordHash   string     // bcrypt — only superadmin/nutritionist; empty for clients
	FirstName      string
	LastName       string
	Gender         string
	BirthDate      *time.Time
	Height         *float64   // cm
	Weight         *float64   // kg
	AvatarURL      string
	IsActive       bool
	NutritionistID *uuid.UUID // set for clients; nil for nutritionists/superadmin
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// BMI calculates body mass index. Returns 0 if height or weight is not set.
func (u *User) BMI() float64 {
	if u.Height == nil || u.Weight == nil || *u.Height == 0 {
		return 0
	}
	heightM := *u.Height / 100
	return *u.Weight / (heightM * heightM)
}

// IsClient returns true if the user has the client role.
func (u *User) IsClient() bool {
	return u.Role == RoleClient
}

// IsNutritionist returns true if the user has the nutritionist role.
func (u *User) IsNutritionist() bool {
	return u.Role == RoleNutritionist
}

// IsSuperAdmin returns true if the user has the superadmin role.
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// BelongsTo returns true if the client belongs to the given nutritionist.
func (u *User) BelongsTo(nutritionistID uuid.UUID) bool {
	if u.NutritionistID == nil {
		return false
	}
	return *u.NutritionistID == nutritionistID
}
