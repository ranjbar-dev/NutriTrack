package model

// UserRole represents the role of a user in the system.
type UserRole string

const (
	RoleSuperAdmin   UserRole = "super_admin"
	RoleNutritionist UserRole = "nutritionist"
	RoleClient       UserRole = "client"
)

// IsValid checks if the role is a known valid role.
func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleNutritionist, RoleClient:
		return true
	}
	return false
}

// GenderType represents the biological gender of a user.
type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
)

// IsValid checks if the gender is a known valid gender.
func (g GenderType) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale:
		return true
	}
	return false
}

// TokenPair holds an access and refresh token pair returned after authentication.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
