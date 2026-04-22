package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/interfaces/http/dto"
	"github.com/rs/zerolog/log"
)

// ErrorHandler centralizes all error responses.
// Handlers call c.Error(err) or c.AbortWithError(), then this middleware converts them.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		requestID, _ := c.Get(RequestIDKey)

		var appErr *shared.AppError
		if errors.As(err, &appErr) {
			log.Debug().
				Str("request_id", requestID.(string)).
				Str("code", appErr.Code).
				Msg("app error")
			dto.Error(c, appErr)
			return
		}

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "اطلاعات ورودی نامعتبر است",
				"errors":  formatValidationErrors(validationErrs),
			})
			return
		}

		log.Error().
			Str("request_id", requestID.(string)).
			Err(err).
			Msg("unhandled error")
		dto.Error(c, shared.ErrInternal)
	}
}

func formatValidationErrors(errs validator.ValidationErrors) []map[string]string {
	out := make([]map[string]string, 0, len(errs))
	for _, e := range errs {
		out = append(out, map[string]string{
			"field":   e.Field(),
			"tag":     e.Tag(),
			"message": persianValidationMessage(e),
		})
	}
	return out
}

func persianValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "این فیلد الزامی است"
	case "email":
		return "فرمت ایمیل نامعتبر است"
	case "min":
		return "مقدار وارد شده کمتر از حداقل مجاز است"
	case "max":
		return "مقدار وارد شده بیشتر از حداکثر مجاز است"
	case "len":
		return "طول مقدار وارد شده نامعتبر است"
	case "numeric":
		return "مقدار وارد شده باید عدد باشد"
	case "gte":
		return "مقدار وارد شده باید بزرگتر یا مساوی حداقل باشد"
	case "lte":
		return "مقدار وارد شده باید کوچکتر یا مساوی حداکثر باشد"
	case "oneof":
		return "مقدار وارد شده در لیست مقادیر مجاز نیست"
	case "url":
		return "فرمت آدرس اینترنتی نامعتبر است"
	case "uuid4":
		return "فرمت شناسه (UUID) نامعتبر است"
	default:
		return "مقدار وارد شده نامعتبر است"
	}
}
