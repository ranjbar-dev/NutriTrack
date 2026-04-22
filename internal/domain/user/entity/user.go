package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Role is a domain type for user roles.
type Role string

// Role constants
const (
	RoleSuperAdmin   Role = "superadmin"
	RoleNutritionist Role = "nutritionist"
	RoleClient       Role = "client"
)

// Gender constants
const (
	GenderMale   = "male"
	GenderFemale = "female"
)

// ErrInvalidRole is returned when an unrecognized role is provided.
var ErrInvalidRole = errors.New("invalid user role")

// ErrMobileRequired is returned when mobile is empty in NewUser.
var ErrMobileRequired = errors.New("mobile is required")

// User is the root aggregate for identity and profile data.
// All fields are unexported; access is via getter/setter methods.
type User struct {
	id             uuid.UUID
	role           Role
	mobile         string     // Iranian mobile, E.164 format: +989xxxxxxxxx
	email          string     // nullable — only superadmin/nutritionist
	passwordHash   string     // bcrypt — only superadmin/nutritionist; empty for clients
	firstName      string
	lastName       string
	gender         string
	birthDate      *time.Time
	height         *float64   // cm
	weight         *float64   // kg
	avatarURL      string
	isActive       bool
	nutritionistID *uuid.UUID // set for clients; nil for nutritionists/superadmin
	createdAt      time.Time
	updatedAt      time.Time
}

// NewUser creates a validated User aggregate. Used by the application layer.
func NewUser(role Role, mobile, firstName, lastName string) (*User, error) {
	if mobile == "" {
		return nil, ErrMobileRequired
	}
	switch role {
	case RoleSuperAdmin, RoleNutritionist, RoleClient:
	default:
		return nil, ErrInvalidRole
	}
	return &User{
		id:        uuid.New(),
		role:      role,
		mobile:    mobile,
		firstName: firstName,
		lastName:  lastName,
		isActive:  true,
	}, nil
}

// Reconstitute rebuilds a User from persisted data. Used only by the repository layer.
func Reconstitute(
	id uuid.UUID, role Role, mobile, email, passwordHash string,
	firstName, lastName, gender string,
	birthDate *time.Time,
	height, weight *float64,
	avatarURL string,
	isActive bool,
	nutritionistID *uuid.UUID,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:             id,
		role:           role,
		mobile:         mobile,
		email:          email,
		passwordHash:   passwordHash,
		firstName:      firstName,
		lastName:       lastName,
		gender:         gender,
		birthDate:      birthDate,
		height:         height,
		weight:         weight,
		avatarURL:      avatarURL,
		isActive:       isActive,
		nutritionistID: nutritionistID,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// --- Getters ---

func (u *User) GetID() uuid.UUID            { return u.id }
func (u *User) GetRole() Role               { return u.role }
func (u *User) GetMobile() string           { return u.mobile }
func (u *User) GetEmail() string            { return u.email }
func (u *User) GetPasswordHash() string     { return u.passwordHash }
func (u *User) GetFirstName() string        { return u.firstName }
func (u *User) GetLastName() string         { return u.lastName }
func (u *User) GetGender() string           { return u.gender }
func (u *User) GetBirthDate() *time.Time    { return u.birthDate }
func (u *User) GetHeight() *float64         { return u.height }
func (u *User) GetWeight() *float64         { return u.weight }
func (u *User) GetAvatarURL() string        { return u.avatarURL }
func (u *User) GetIsActive() bool           { return u.isActive }
func (u *User) GetNutritionistID() *uuid.UUID { return u.nutritionistID }
func (u *User) GetCreatedAt() time.Time     { return u.createdAt }
func (u *User) GetUpdatedAt() time.Time     { return u.updatedAt }

// --- Domain read methods ---

// FullName returns the user's full name.
func (u *User) FullName() string {
	return u.firstName + " " + u.lastName
}

// BMI calculates body mass index. Returns 0 if height or weight is not set.
func (u *User) BMI() float64 {
	if u.height == nil || u.weight == nil || *u.height == 0 {
		return 0
	}
	heightM := *u.height / 100
	return *u.weight / (heightM * heightM)
}

// IsClient returns true if the user has the client role.
func (u *User) IsClient() bool {
	return u.role == RoleClient
}

// IsNutritionist returns true if the user has the nutritionist role.
func (u *User) IsNutritionist() bool {
	return u.role == RoleNutritionist
}

// IsSuperAdmin returns true if the user has the superadmin role.
func (u *User) IsSuperAdmin() bool {
	return u.role == RoleSuperAdmin
}

// BelongsTo returns true if the client belongs to the given nutritionist.
func (u *User) BelongsTo(nutritionistID uuid.UUID) bool {
	if u.nutritionistID == nil {
		return false
	}
	return *u.nutritionistID == nutritionistID
}

// --- Domain mutation methods ---

// Activate enables the user account.
func (u *User) Activate() { u.isActive = true }

// Deactivate disables the user account.
func (u *User) Deactivate() { u.isActive = false }

// SetActive sets the active status directly (useful for conditional toggling).
func (u *User) SetActive(active bool) { u.isActive = active }

// SetPasswordHash updates the stored bcrypt hash.
func (u *User) SetPasswordHash(hash string) { u.passwordHash = hash }

// SetAvatarURL updates the stored avatar URL.
func (u *User) SetAvatarURL(url string) { u.avatarURL = url }

// SetEmail updates the user's email address.
func (u *User) SetEmail(email string) { u.email = email }

// AssignNutritionist links a client to a nutritionist.
func (u *User) AssignNutritionist(id uuid.UUID) { u.nutritionistID = &id }

// UpdateProfile applies partial profile updates; empty/nil values are ignored.
func (u *User) UpdateProfile(firstName, lastName, gender string, birthDate *time.Time, height, weight *float64) {
	if firstName != "" {
		u.firstName = firstName
	}
	if lastName != "" {
		u.lastName = lastName
	}
	if gender != "" {
		u.gender = gender
	}
	if birthDate != nil {
		u.birthDate = birthDate
	}
	if height != nil {
		u.height = height
	}
	if weight != nil {
		u.weight = weight
	}
}

// --- Infrastructure-only setters (for repository to populate DB-generated values) ---

// SetID sets the aggregate ID (used by the repository after DB insert).
func (u *User) SetID(id uuid.UUID) { u.id = id }

// SetCreatedAt sets the creation timestamp (used by the repository after DB insert).
func (u *User) SetCreatedAt(t time.Time) { u.createdAt = t }

// SetUpdatedAt sets the updated timestamp (used by the repository after DB insert/update).
func (u *User) SetUpdatedAt(t time.Time) { u.updatedAt = t }
