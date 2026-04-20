package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetAuthCookiesUsesStrictSameSite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	setAuthCookies(ctx, "access-token", "refresh-token")

	cookies := rec.Header().Values("Set-Cookie")
	require.Len(t, cookies, 2)
	require.Contains(t, strings.Join(cookies, "\n"), "access_token=access-token")
	require.Contains(t, strings.Join(cookies, "\n"), "refresh_token=refresh-token")
	require.Contains(t, strings.Join(cookies, "\n"), "Path=/api")
	require.Contains(t, strings.Join(cookies, "\n"), "Path=/api/auth/refresh")
	require.Contains(t, strings.Join(cookies, "\n"), "SameSite=Strict")
}

func TestClearAuthCookiesKeepsStrictSameSite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	clearAuthCookies(ctx)

	cookies := rec.Header().Values("Set-Cookie")
	require.Len(t, cookies, 2)
	require.Contains(t, strings.Join(cookies, "\n"), "Max-Age=0")
	require.Contains(t, strings.Join(cookies, "\n"), "SameSite=Strict")
}
