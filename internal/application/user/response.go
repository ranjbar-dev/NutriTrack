package user

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
)

// MapClientResponse converts a domain User (client role) to a JSON-serialisable map.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapClientResponse(u *entity.User) map[string]any {
	var birthDate *string
	if u.GetBirthDate() != nil {
		s := u.GetBirthDate().Format("2006-01-02")
		birthDate = &s
	}

	var nutID *string
	if u.GetNutritionistID() != nil {
		s := u.GetNutritionistID().String()
		nutID = &s
	}

	return map[string]any{
		"id":              u.GetID(),
		"mobile":          u.GetMobile(),
		"first_name":      u.GetFirstName(),
		"last_name":       u.GetLastName(),
		"full_name":       u.FullName(),
		"gender":          u.GetGender(),
		"birth_date":      birthDate,
		"height":          u.GetHeight(),
		"weight":          u.GetWeight(),
		"bmi":             u.BMI(),
		"avatar_url":      u.GetAvatarURL(),
		"is_active":       u.GetIsActive(),
		"nutritionist_id": nutID,
		"created_at":      u.GetCreatedAt(),
		"updated_at":      u.GetUpdatedAt(),
	}
}

// MapNutritionistResponse converts a domain User (nutritionist role) to a JSON-serialisable map.
func MapNutritionistResponse(u *entity.User) map[string]any {
	return map[string]any{
		"id":         u.GetID(),
		"email":      u.GetEmail(),
		"mobile":     u.GetMobile(),
		"first_name": u.GetFirstName(),
		"last_name":  u.GetLastName(),
		"is_active":  u.GetIsActive(),
		"created_at": u.GetCreatedAt(),
	}
}
