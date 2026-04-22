# DDD Fix Report: internal/infrastructure/sms
Layer: infrastructure
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | `NewKavenegarAdapter` omits required-field validation and error return | MEDIUM | kavenegar_adapter.go, bootstrap/wire.go | SAFE | FIXED |

## Changes Applied

### Fix 1: Add required-field validation and error return to `NewKavenegarAdapter`

**File:** `internal/infrastructure/sms/kavenegar_adapter.go`
**Change:** Added `strings` import; changed return type from `*KavenegarAdapter` to `(*KavenegarAdapter, error)`; added `strings.TrimSpace` validation for both `apiKey` and `otpTemplate` before constructing the struct.

**Before:**
```go
func NewKavenegarAdapter(apiKey, otpTemplate string) *KavenegarAdapter {
    return &KavenegarAdapter{
        apiKey:      apiKey,
        otpTemplate: otpTemplate,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}
```

**After:**
```go
func NewKavenegarAdapter(apiKey, otpTemplate string) (*KavenegarAdapter, error) {
    if strings.TrimSpace(apiKey) == "" {
        return nil, fmt.Errorf("kavenegar: apiKey is required")
    }
    if strings.TrimSpace(otpTemplate) == "" {
        return nil, fmt.Errorf("kavenegar: otpTemplate is required")
    }
    return &KavenegarAdapter{
        apiKey:      apiKey,
        otpTemplate: otpTemplate,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }, nil
}
```

**Build:** PASS

---

**File:** `bootstrap/wire.go`
**Change:** Added `"github.com/rs/zerolog/log"` import; updated the Kavenegar construction site to handle the returned error via `log.Fatal`.

**Before:**
```go
if cfg.App.Env == "production" && cfg.SMS.KavenegarAPIKey != "" {
    smsProvider = sms.NewKavenegarAdapter(cfg.SMS.KavenegarAPIKey, cfg.SMS.OTPTemplate)
}
```

**After:**
```go
if cfg.App.Env == "production" && cfg.SMS.KavenegarAPIKey != "" {
    kavAdapter, err := sms.NewKavenegarAdapter(cfg.SMS.KavenegarAPIKey, cfg.SMS.OTPTemplate)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to initialize Kavenegar SMS adapter")
    }
    smsProvider = kavAdapter
}
```

**Build:** PASS

## Deferred Items
None.

## Final Build Status
PASS — `go build ./...` after all fixes
N/A — `go vet` not explicitly run (build is clean)

## Remaining Violations
None — the single MEDIUM finding is resolved. DDD-AUDIT.md may be removed.
