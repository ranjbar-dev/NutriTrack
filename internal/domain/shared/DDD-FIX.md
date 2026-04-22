# DDD Fix Report: internal/domain/shared
Layer: domain + interface (callers updated)
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes (exit code 0)

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | `AppError` imports `net/http` + carries `HTTPStatus` | CRITICAL | `apperror.go`, `dto/response.go`, `middleware/*.go` | SAFE | FIXED |
| 2 | `json:` struct tags on `AppError` | CRITICAL | `apperror.go` | SAFE | FIXED |
| 3 | `ToResponse()` HTTP serialization logic in domain | HIGH | `apperror.go`, `dto/response.go`, `middleware/*.go` | SAFE | FIXED |
| 4 | `ImageInfo` value object has exported (mutable) fields | HIGH | `file_validator.go`, `avatar_service.go` | SAFE | FIXED |
| 5 | `panic` in `init()` for timezone loading | LOW | (not in audit scope for this pass) | DEFERRED | DEFERRED |

## Changes Applied

### Fix 1 & 2 & 3: Remove `net/http`, `HTTPStatus`, `json:` tags, and `ToResponse()` from `AppError`

**File:** `internal/domain/shared/apperror.go`

**Before:**
```go
import "net/http"

type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"` // Always in Persian
    HTTPStatus int    `json:"-"`
}

func (e *AppError) ToResponse() map[string]any {
    return map[string]any{"code": e.Code, "message": e.Message}
}

var ErrInternal = &AppError{
    Code:       "INTERNAL_ERROR",
    Message:    "خطای داخلی سرور رخ داده است",
    HTTPStatus: http.StatusInternalServerError,
}
// ... all other vars had HTTPStatus: http.StatusXxx

func (e *AppError) WithMessage(msg string) *AppError {
    return &AppError{Code: e.Code, Message: msg, HTTPStatus: e.HTTPStatus}
}
```

**After:**
```go
// no import

type AppError struct {
    Code    string
    Message string // Always in Persian
}

// ToResponse() deleted

var ErrInternal = &AppError{
    Code:    "INTERNAL_ERROR",
    Message: "خطای داخلی سرور رخ داده است",
}
// ... all other vars have no HTTPStatus field

func (e *AppError) WithMessage(msg string) *AppError {
    return &AppError{Code: e.Code, Message: msg}
}
```

**Build after this change:** PASS (after updating callers below)

---

### Fix 1 & 3 (callers): Add `httpStatusFor()` mapping + update `Error()`/`Abort()` in interface layer

**File:** `internal/interfaces/http/dto/response.go`

**Change:** Added `httpStatusFor(*shared.AppError) int` function that maps domain error codes to HTTP status codes. Updated `Error()` and `Abort()` to use it and inline the JSON body instead of calling `ToResponse()`.

```go
// NEW — HTTP status mapping lives entirely in the interface layer
func httpStatusFor(err *shared.AppError) int {
    switch err.Code {
    case "NOT_FOUND", "USER_NOT_FOUND", ...:
        return http.StatusNotFound
    case "UNAUTHORIZED", "INVALID_CREDENTIALS", ...:
        return http.StatusUnauthorized
    // ... etc
    default:
        return http.StatusInternalServerError
    }
}

func Error(c *gin.Context, err *shared.AppError) {
    c.JSON(httpStatusFor(err), gin.H{"code": err.Code, "message": err.Message})
}

func Abort(c *gin.Context, err *shared.AppError) {
    c.AbortWithStatusJSON(httpStatusFor(err), gin.H{"code": err.Code, "message": err.Message})
}
```

**Build:** PASS

---

### Fix 1 & 3 (callers): Update `middleware/error_handler.go`

**File:** `internal/interfaces/http/middleware/error_handler.go`

**Change:** Added `dto` import. Replaced `c.JSON(appErr.HTTPStatus, appErr.ToResponse())` with `dto.Error(c, appErr)`. Replaced `c.JSON(http.StatusInternalServerError, shared.ErrInternal.ToResponse())` with `dto.Error(c, shared.ErrInternal)`.

**Build:** PASS

---

### Fix 1 & 3 (callers): Update `middleware/recovery.go`

**File:** `internal/interfaces/http/middleware/recovery.go`

**Change:** Added `dto` import, removed `net/http` import. Replaced `c.JSON(http.StatusInternalServerError, shared.ErrInternal.ToResponse())` with `dto.Error(c, shared.ErrInternal)`.

**Build:** PASS

---

### Fix 1 & 3 (callers): Update `middleware/not_found.go`

**File:** `internal/interfaces/http/middleware/not_found.go`

**Change:** Added `dto` import, removed `net/http` import. Replaced `c.JSON(http.StatusNotFound, shared.ErrNotFound.ToResponse())` with `dto.Error(c, shared.ErrNotFound)`.

**Build:** PASS

---

### Fix 1 (caller): Remove `HTTPStatus` from `ErrRateLimitExceeded` in `middleware/rate_limit.go`

**File:** `internal/interfaces/http/middleware/rate_limit.go`

**Before:**
```go
var ErrRateLimitExceeded = &shared.AppError{
    Code:       "RATE_LIMIT_EXCEEDED",
    Message:    "تعداد درخواست‌های شما بیش از حد مجاز است",
    HTTPStatus: 429,
}
```

**After:**
```go
var ErrRateLimitExceeded = &shared.AppError{
    Code:    "RATE_LIMIT_EXCEEDED",
    Message: "تعداد درخواست‌های شما بیش از حد مجاز است",
}
```

Note: `RATE_LIMIT_EXCEEDED` is mapped to `http.StatusTooManyRequests` (429) in `httpStatusFor()`.

**Build:** PASS

---

### Fix 4: Make `ImageInfo` fields unexported; add `NewImageInfo()` constructor and getters

**File:** `internal/domain/shared/file_validator.go`

**Before:**
```go
type ImageInfo struct {
    Extension string // "jpg", "png", "webp"
    MIMEType  string // "image/jpeg", "image/png", "image/webp"
}

// ValidateImageMagicBytes returned &ImageInfo{Extension: "jpg", MIMEType: "image/jpeg"}
```

**After:**
```go
type ImageInfo struct {
    extension string // "jpg", "png", "webp"
    mimeType  string // "image/jpeg", "image/png", "image/webp"
}

func NewImageInfo(extension, mimeType string) *ImageInfo {
    return &ImageInfo{extension: extension, mimeType: mimeType}
}

func (i *ImageInfo) Extension() string { return i.extension }
func (i *ImageInfo) MIMEType() string  { return i.mimeType }

// ValidateImageMagicBytes now returns NewImageInfo("jpg", "image/jpeg")
```

**Build:** PASS

---

### Fix 4 (caller): Update `avatar_service.go` to use getter

**File:** `internal/application/user/avatar_service.go`

**Before:** `s.storage.SaveAvatar(reader, imgInfo.Extension)`
**After:** `s.storage.SaveAvatar(reader, imgInfo.Extension())`

**Build:** PASS

---

## Deferred Items

| # | Finding | Reason |
|---|---------|--------|
| 5 | `panic` in `init()` for timezone loading | LOW priority; not in scope for CRITICAL/HIGH pass |

## Final Build Status
PASS — `go build ./...` after all fixes (exit code 0)
PASS — `go vet ./internal/...` after all fixes (exit code 0)

## Remaining Violations
None — all CRITICAL and HIGH findings have been resolved.
