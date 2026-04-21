package shared

import "time"

var tehranLoc *time.Location

func init() {
	var err error
	tehranLoc, err = time.LoadLocation("Asia/Tehran")
	if err != nil {
		// Should never happen since time/tzdata is embedded in main.go
		panic("failed to load Asia/Tehran timezone: " + err.Error())
	}
}

// NowTehran returns current time in Asia/Tehran timezone.
func NowTehran() time.Time {
	return time.Now().In(tehranLoc)
}

// TodayTehran returns today's date in Asia/Tehran timezone.
// IMPORTANT: Do NOT use time.Now().UTC() for date logic — after 20:30 UTC
// it returns tomorrow's date in Iran. Always use this function for "today" in Iran.
func TodayTehran() time.Time {
	now := NowTehran()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tehranLoc)
}

// TehranLoc returns the Asia/Tehran *time.Location for use in time formatting.
func TehranLoc() *time.Location {
	return tehranLoc
}

// ToTehran converts any time.Time to Asia/Tehran timezone.
func ToTehran(t time.Time) time.Time {
	return t.In(tehranLoc)
}
