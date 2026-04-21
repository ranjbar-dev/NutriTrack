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

// Error sends an AppError response.
func Error(c *gin.Context, err *shared.AppError) {
	c.JSON(err.HTTPStatus, err.ToResponse())
}

// Abort aborts with an AppError response.
func Abort(c *gin.Context, err *shared.AppError) {
	c.AbortWithStatusJSON(err.HTTPStatus, err.ToResponse())
}
