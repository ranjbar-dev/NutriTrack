package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// OK sends a 200 response with data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

// NoContent sends a 204 response.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Paginated sends a 200 response with data + pagination metadata.
func Paginated(c *gin.Context, data any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// httpStatusFor maps a domain AppError code to its HTTP transport status code.
// This is the interface layer's authoritative translation of domain errors to HTTP semantics.
func httpStatusFor(err *shared.AppError) int {
	switch err.Code {
	case "NOT_FOUND", "USER_NOT_FOUND", "FOOD_NOT_FOUND", "MEDICATION_NOT_FOUND",
		"DIET_PLAN_NOT_FOUND", "TRACKING_NOT_FOUND", "LAB_RESULT_NOT_FOUND",
		"MESSAGE_NOT_FOUND", "FOOD_REQUEST_NOT_FOUND", "NOTIFICATION_PREFERENCE_NOT_FOUND":
		return http.StatusNotFound
	case "UNAUTHORIZED", "INVALID_CREDENTIALS", "INVALID_TOKEN", "TOKEN_REVOKED",
		"OTP_INVALID", "OTP_EXPIRED", "OTP_MAX_ATTEMPTS":
		return http.StatusUnauthorized
	case "FORBIDDEN", "FOOD_REQUEST_NOT_OWNED":
		return http.StatusForbidden
	case "CONFLICT", "USER_ALREADY_EXISTS", "FOOD_REQUEST_ALREADY_PROCESSED", "PLAN_ALREADY_ACTIVE":
		return http.StatusConflict
	case "VALIDATION_ERROR", "INVALID_MOBILE", "INVALID_FILE_TYPE":
		return http.StatusUnprocessableEntity
	case "OTP_RATE_LIMIT", "RATE_LIMIT_EXCEEDED":
		return http.StatusTooManyRequests
	case "FILE_TOO_LARGE":
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

// Error sends an AppError response.
func Error(c *gin.Context, err *shared.AppError) {
	c.JSON(httpStatusFor(err), gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}

// Abort aborts with an AppError response.
func Abort(c *gin.Context, err *shared.AppError) {
	c.AbortWithStatusJSON(httpStatusFor(err), gin.H{
		"code":    err.Code,
		"message": err.Message,
	})
}
