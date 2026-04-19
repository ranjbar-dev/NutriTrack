package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// iranianMobileRegex matches Iranian mobile numbers: 09XXXXXXXXX (11 digits).
var iranianMobileRegex = regexp.MustCompile(`^09[0-9]{9}$`)

// RegisterCustomValidators registers all custom validation rules with Gin's validator.
func RegisterCustomValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("iranian_mobile", validateIranianMobile); err != nil {
			return err
		}
	}
	return nil
}

// validateIranianMobile validates that the field value is a valid Iranian mobile number.
func validateIranianMobile(fl validator.FieldLevel) bool {
	return iranianMobileRegex.MatchString(fl.Field().String())
}
