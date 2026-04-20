package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSecurityHeadersSetCommonHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("X-Forwarded-Proto", "https")

	called := false
	SecurityHeaders()(ctx)
	ctx.Next()
	called = true

	require.True(t, called)
	require.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	require.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	require.Equal(t, "default-src 'self'", rec.Header().Get("Content-Security-Policy"))
	require.Equal(t, "same-site", rec.Header().Get("Cross-Origin-Resource-Policy"))
	require.Equal(t, "max-age=63072000; includeSubDomains; preload", rec.Header().Get("Strict-Transport-Security"))
}

func TestSecurityHeadersSkipsHSTSWithoutHTTPS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	SecurityHeaders()(ctx)

	require.Empty(t, rec.Header().Get("Strict-Transport-Security"))
}
