package config

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	appTZOnce sync.Once
	appTZLoc  *time.Location
)

func EnvOrDefault(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func SupportAppLocation() *time.Location {
	appTZOnce.Do(func() {
		tzName := EnvOrDefault("APP_TIMEZONE", EnvOrDefault("read_db_timezone", "Asia/Jakarta"))
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			log.Printf("[WARNING] Invalid timezone %q, fallback to Asia/Jakarta: %v", tzName, err)
			loc = time.FixedZone("Asia/Jakarta", 7*60*60)
		}
		appTZLoc = loc
	})

	if appTZLoc == nil {
		return time.FixedZone("Asia/Jakarta", 7*60*60)
	}

	return appTZLoc
}

func NormalizeNaiveTimestampToAppTZ(ts time.Time) time.Time {
	if ts.IsZero() {
		return ts
	}

	loc := SupportAppLocation()
	year, month, day := ts.Date()
	hour, minute, second := ts.Clock()

	// Re-attach app timezone to DB timestamp fields without shifting wall-clock time.
	return time.Date(year, month, day, hour, minute, second, ts.Nanosecond(), loc)
}
