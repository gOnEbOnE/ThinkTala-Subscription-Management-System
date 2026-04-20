package dashboard

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func parsePeriod(c *gin.Context) (string, time.Time, time.Time, error) {
	now := time.Now()
	loc := now.Location()

	monthRaw := strings.TrimSpace(c.Query("month"))
	yearRaw := strings.TrimSpace(c.Query("year"))
	startRaw := strings.TrimSpace(c.Query("start_date"))
	endRaw := strings.TrimSpace(c.Query("end_date"))

	if monthRaw != "" && yearRaw == "" {
		return "", time.Time{}, time.Time{}, errors.New("month wajib dengan year")
	}

	if startRaw != "" || endRaw != "" {
		if startRaw == "" || endRaw == "" {
			return "", time.Time{}, time.Time{}, errors.New("start_date dan end_date wajib diisi bersamaan")
		}

		startDate, err := time.ParseInLocation("2006-01-02", startRaw, loc)
		if err != nil || startDate.Format("2006-01-02") != startRaw {
			return "", time.Time{}, time.Time{}, errors.New("format start_date tidak valid (YYYY-MM-DD)")
		}

		endDate, err := time.ParseInLocation("2006-01-02", endRaw, loc)
		if err != nil || endDate.Format("2006-01-02") != endRaw {
			return "", time.Time{}, time.Time{}, errors.New("format end_date tidak valid (YYYY-MM-DD)")
		}

		if endDate.Before(startDate) {
			return "", time.Time{}, time.Time{}, errors.New("Rentang tanggal tidak valid.")
		}

		return "custom", startDate, endDate.Add(24*time.Hour - time.Nanosecond), nil
	}

	if monthRaw != "" || yearRaw != "" {
		year, err := strconv.Atoi(yearRaw)
		if err != nil || year < 1900 || year > 3000 {
			return "", time.Time{}, time.Time{}, errors.New("year tidak valid")
		}

		if monthRaw == "" {
			start := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
			end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
			return "yearly", start, end, nil
		}

		month, convErr := strconv.Atoi(monthRaw)
		if convErr != nil || month < 1 || month > 12 {
			return "", time.Time{}, time.Time{}, errors.New("month tidak valid")
		}

		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return "monthly", start, end, nil
	}

	period := strings.ToLower(strings.TrimSpace(c.Query("period")))
	if period == "yearly" {
		start := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, loc)
		return "yearly", start, start.AddDate(1, 0, 0).Add(-time.Nanosecond), nil
	}

	base := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	return "monthly", base, base.AddDate(0, 1, 0).Add(-time.Nanosecond), nil
}

func previousPeriod(start, end time.Time) (time.Time, time.Time) {
	duration := end.Sub(start) + time.Nanosecond
	prevEnd := start.Add(-time.Nanosecond)
	prevStart := prevEnd.Add(-duration + time.Nanosecond)
	return prevStart, prevEnd
}

func buildMonthBuckets(startDate, endDate time.Time) []time.Time {
	cursor := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	last := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())
	months := make([]time.Time, 0)
	for !cursor.After(last) {
		months = append(months, cursor)
		cursor = cursor.AddDate(0, 1, 0)
	}
	if months == nil {
		return []time.Time{}
	}
	return months
}

func monthKey(ts time.Time) string {
	return fmt.Sprintf("%04d-%02d", ts.Year(), ts.Month())
}

func calculateRate(numerator, denominator int) float64 {
	if denominator <= 0 {
		return 0
	}
	return round2((float64(numerator) / float64(denominator)) * 100)
}

func compareRate(current, previous int) float64 {
	if previous <= 0 {
		if current <= 0 {
			return 0
		}
		return 100
	}
	return round2(((float64(current) - float64(previous)) / float64(previous)) * 100)
}

func compareFloatRate(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return round2(((current - previous) / previous) * 100)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func calcTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func parseIntDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func normalizePageLimit(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return page, limit
}

func packagePeriodLabel(periodType string) string {
	switch strings.ToLower(strings.TrimSpace(periodType)) {
	case "yearly":
		return "This Year"
	case "custom":
		return "Custom Period"
	default:
		return "This Month"
	}
}

func classifyPackageBucket(packageID, packageName string) string {
	raw := strings.ToLower(strings.TrimSpace(packageID + " " + packageName))
	switch {
	case strings.Contains(raw, "enterprise"):
		return "enterprise"
	case strings.Contains(raw, "premium"), strings.Contains(raw, "pro"):
		return "premium"
	case strings.Contains(raw, "starter"), strings.Contains(raw, "basic"), strings.Contains(raw, "trial"):
		return "starter"
	default:
		return "starter"
	}
}
