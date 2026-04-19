---
phase: 01-foundation-infrastructure
plan: 03
title: "Go Auth Packages & Middleware Suite"
subsystem: backend
tags: [go, gin, jwt, sms, middleware, auth, rate-limiting, cors, security]
dependency_graph:
  requires: [01-01]
  provides: [jwt-package, sms-sender-interface, auth-middleware, role-guard, rate-limiter, cors-middleware, request-id-middleware, logger-middleware, security-headers]
  affects: [01-04, 01-05]
tech_stack:
  added: [golang-jwt/v5@5.3.1, google/uuid@1.6.0]
  patterns: [httpOnly-cookie-auth, sliding-window-rate-limit, phone-normalization, structured-json-logging, middleware-chain]
key_files:
  created:
    - backend/pkg/jwt/jwt.go
    - backend/pkg/sms/sms.go
    - backend/pkg/sms/mock.go
    - backend/pkg/sms/kavenegar.go
    - backend/internal/middleware/auth.go
    - backend/internal/middleware/role_guard.go
    - backend/internal/middleware/rate_limit.go
    - backend/internal/middleware/cors.go
    - backend/internal/middleware/request_id.go
    - backend/internal/middleware/logger.go
    - backend/internal/middleware/security_headers.go
  modified:
    - backend/go.mod
    - backend/go.sum
decisions:
  - "JWT tokens signed with HMAC-SHA256 only; ParseToken explicitly validates signing method to prevent algorithm confusion"
  - "Rate limiter reads mobile from JSON body (peeks without consuming) then normalizes phone before keying"
  - "Security headers middleware added (Rule 2) since critical constraints require X-Content-Type-Options, X-Frame-Options, etc."
metrics:
  duration: "4 minutes"
  completed: "2026-04-19T15:54:00Z"
---

# Phase 01 Plan 03: Go Auth Packages & Middleware Suite Summary

JWT utility with HMAC-SHA256 access/refresh tokens, SMS sender interface with mock/Kavenegar adapters, and 7-middleware Gin suite (auth, role guard, rate limiter, CORS, request ID, logger, security headers).

## What Was Done

### Task 1: JWT utility package and SMS sender abstraction
- Created `pkg/jwt/jwt.go` with `Claims` struct extending `jwt.RegisteredClaims` (user_id, role)
- `CreateAccessToken` signs with HS256, 15-minute expiry, issuer "nutritrack"
- `CreateRefreshToken` signs with HS256, 30-day expiry
- `CreateTokenPair` creates both and returns `model.TokenPair`
- `ParseToken` validates HMAC signing method explicitly (prevents algorithm confusion attack T-03-01)
- Created `pkg/sms/sms.go` with `Sender` interface (`SendOTP(phone, code string) error`)
- Created `pkg/sms/mock.go` with `MockSender` that logs OTP to stdout via zerolog (dev mode per D-04)
- Created `pkg/sms/kavenegar.go` with `KavenegarSender` using Kavenegar verify/lookup REST API (10s timeout)
- **Commit:** `d8eb046`

### Task 2: Auth and role guard middleware
- Created `internal/middleware/auth.go` — extracts JWT from "access_token" httpOnly cookie (D-01)
- Sets "user_id" and "role" in Gin context for downstream handlers
- Returns 401 with Persian error ("احراز هویت الزامی است") for missing cookies
- Returns 401 with Persian error ("توکن نامعتبر است") for invalid/expired tokens
- Created `internal/middleware/role_guard.go` — checks role against allowed roles list
- Returns 403 with Persian error ("دسترسی غیرمجاز") if role doesn't match
- **Commit:** `c0f80a5`

### Task 3: Rate limiter, CORS, request ID, logger, and security headers middleware
- Created `internal/middleware/rate_limit.go`:
  - `RateLimiter` struct with sync.Mutex, sliding window entries map
  - `NewRateLimiter(max, window)` with background cleanup goroutine (every 1 min)
  - `Allow(key)` checks sliding window count against max
  - `RateLimit()` middleware reads "mobile" from JSON body, normalizes phone, keys on it
  - Phone normalization: strips +98/0098/0 prefix → canonical 10-digit 9XXXXXXXXX format (T-03-02, Pitfall 3)
  - Falls back to client IP if no mobile in body
  - Returns 429 with Persian error for exceeded limit
- Created `internal/middleware/cors.go`:
  - Sets explicit frontend URL origin (not "*") with credentials=true (T-03-03, SEC-06)
  - Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
  - Allow-Headers: Content-Type, Authorization
  - Max-Age: 86400 seconds
  - Handles OPTIONS preflight with 204 No Content
- Created `internal/middleware/request_id.go`:
  - Generates UUID v4 via google/uuid
  - Sets "request_id" in Gin context and X-Request-ID response header (T-03-06)
- Created `internal/middleware/logger.go`:
  - Structured JSON logging via zerolog (D-23, INFRA-04)
  - Fields: request_id, method, path, status, duration_ms, client_ip
  - Timestamp added automatically by zerolog
- Created `internal/middleware/security_headers.go` (Rule 2 deviation):
  - X-Content-Type-Options: nosniff
  - X-Frame-Options: DENY
  - X-XSS-Protection: 1; mode=block
  - Referrer-Policy: strict-origin-when-cross-origin
  - Content-Security-Policy: default-src 'self'
- **Commit:** `cd22cb7`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing functionality] Security headers middleware**
- **Found during:** Task 3
- **Issue:** Critical constraints specified "Security headers middleware (X-Content-Type-Options, X-Frame-Options, etc.)" but the plan's files_modified list did not include a security_headers.go file
- **Fix:** Added `internal/middleware/security_headers.go` with standard security headers
- **Files modified:** `backend/internal/middleware/security_headers.go`
- **Commit:** `cd22cb7`

## Verification Results

```
cd backend && go build ./pkg/...              → PASS
cd backend && go vet ./pkg/...                → PASS
cd backend && go build ./internal/middleware/  → PASS
cd backend && go vet ./internal/middleware/    → PASS
cd backend && go build ./...                  → PASS (full project)
cd backend && go vet ./...                    → PASS (full project)
```

## Self-Check: PASSED

All 11 created files verified present. All 3 task commits verified in git log.
