package valueobject

import (
	"regexp"
	"strings"

	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
)

// mobileRegex matches Iranian mobile numbers:
// - Starting with 09 (local) or +989 (E.164)
// - 11 digits local / 13 digits E.164
var mobileRegex = regexp.MustCompile(`^(\+989|09)\d{9}$`)

// Mobile is an immutable value object for Iranian mobile numbers.
type Mobile struct {
	value string // normalized to +989xxxxxxxxx format
}

// NewMobile validates and normalizes an Iranian mobile number.
func NewMobile(raw string) (Mobile, error) {
	raw = strings.TrimSpace(raw)
	if !mobileRegex.MatchString(raw) {
		return Mobile{}, shared.ErrInvalidMobile
	}
	// Normalize: 09xxxxxxxxx → +989xxxxxxxxx
	if strings.HasPrefix(raw, "09") {
		raw = "+98" + raw[1:]
	}
	return Mobile{value: raw}, nil
}

// String returns the normalized E.164 mobile number.
func (m Mobile) String() string {
	return m.value
}

// Local returns the local format (09xxxxxxxxx).
func (m Mobile) Local() string {
	if strings.HasPrefix(m.value, "+98") {
		return "0" + m.value[3:]
	}
	return m.value
}
