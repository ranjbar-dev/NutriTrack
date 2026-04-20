package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthorizationAuditSensitiveRoutesHaveOwnershipGuards(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		snippets []string
	}{
		{
			name: "tracking handler uses nutritionist ownership helper",
			path: filepath.Join("..", "handler", "tracking_handler.go"),
			snippets: []string{
				"nutriClientAndOwner(c)",
				"GetLabResultForNutritionist",
				"ErrTrackingUnauthorized",
			},
		},
		{
			name: "client profile handler delegates ownership-aware service methods",
			path: filepath.Join("..", "handler", "client_handler.go"),
			snippets: []string{
				"GetClientProfile",
				"UpdateClientProfile",
				"ListClients",
			},
		},
		{
			name: "diet plan handler reads auth user and forwards authorization context",
			path: filepath.Join("..", "handler", "diet_plan_handler.go"),
			snippets: []string{
				`c.GetString("role")`,
				`c.GetString("user_id")`,
				"GetPlanAggregate",
			},
		},
		{
			name: "route wiring keeps sensitive groups behind auth and role guards",
			path: filepath.Join("..", "..", "cmd", "api", "main.go"),
			snippets: []string{
				`nutri.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist"))`,
				`dietPlanRead.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin", "client"))`,
				`msgs.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client", "nutritionist"))`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			content, err := os.ReadFile(tc.path)
			require.NoError(t, err)

			text := string(content)
			for _, snippet := range tc.snippets {
				require.Truef(t, strings.Contains(text, snippet), "expected %s to contain %q", tc.path, snippet)
			}
		})
	}
}
