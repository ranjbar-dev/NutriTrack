# 05-04 Summary: User Service & Client Handler Updates

## What Was Built

### User Service (`backend/internal/service/user_service.go`)

Added to existing service:
- `ListClients(nutriID uuid.UUID, params SearchClientsParams)` — delegates to user repo SearchClients
- `GetClientProfile(nutriID, clientID uuid.UUID)` — verifies ownership, returns ClientProfileResponse
- `ActivateClient(nutriID, clientID uuid.UUID)` — sets is_active = true
- `DeactivateClient(nutriID, clientID uuid.UUID)` — sets is_active = false
- `UpdateClientProfile(nutriID, clientID uuid.UUID, req dto.UpdateClientProfileRequest)` — updates height/dob/gender
- `sqlcUserToClientProfile(u sqlc.User)` — converts sqlc User to ClientProfileResponse DTO
- `normalizeUserNotFound(err error)` — maps pgx ErrNoRows → ErrUserNotFound (shared var from auth_service.go)

### Client Handler (`backend/internal/handler/client_handler.go`)

Full rewrite adding 5 new handlers alongside existing `RegisterClient`:
- `NutriListClients` — GET /api/nutritionist/clients with search/filter/sort query params
- `NutriGetClientProfile` — GET /api/nutritionist/clients/:clientId
- `NutriActivateClient` — PUT /api/nutritionist/clients/:clientId/activate
- `NutriDeactivateClient` — PUT /api/nutritionist/clients/:clientId/deactivate
- `NutriUpdateClientProfile` — PUT /api/nutritionist/clients/:clientId/profile

Helper functions:
- `parseNutritionistID(c *gin.Context)` — extracts and parses nutritionist UUID from JWT claims
- `parseClientParam(c *gin.Context)` — parses `:clientId` param as UUID

## Key Technical Details

- `ErrUserNotFound` already declared in `auth_service.go` (same package) — not redeclared in user_service.go
- Handler parses JWT claims directly (role + user_id from middleware-set context keys)
- All route parameters validated as UUID before service calls

## Outcomes

- Client management backend complete
- `go build ./...` passes
