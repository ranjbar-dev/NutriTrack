package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics and returns a 500 error.
// Uses Persian error message consistent with the application's language.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "خطای داخلی سرور",
		})
	})
}
