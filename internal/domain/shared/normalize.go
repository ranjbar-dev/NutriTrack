package shared

import "strings"

// NormalizePersian normalizes Arabic characters to their Persian equivalents.
// This is required for consistent Persian text search with pg_trgm.
// Arabic Kaf (ك) → Persian Kaf (ک)
// Arabic Yeh (ي) → Persian Yeh (ی)
func NormalizePersian(s string) string {
	s = strings.ReplaceAll(s, "ك", "ک")
	s = strings.ReplaceAll(s, "ي", "ی")
	return s
}
