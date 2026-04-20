package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	ctx.Request.Header.Set("Origin", "https://app.example.com")

	CORS("https://app.example.com")(ctx)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "Origin", rec.Header().Get("Vary"))
	require.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSRejectsUnexpectedPreflightOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	ctx.Request.Header.Set("Origin", "https://evil.example.com")

	CORS("https://app.example.com")(ctx)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "مبدا درخواست مجاز نیست")
	require.Equal(t, "Origin", rec.Header().Get("Vary"))
}
