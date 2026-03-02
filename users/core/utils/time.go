package utils

import (
	"fmt"
	"math"
	"time"
)

var monthsID = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// FormatDateID mengubah time.Time ke format "21 Januari 2026"
func FormatDateID(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	// Pastikan menggunakan waktu lokal server
	// t = t.In(time.Local)

	day := t.Day()
	month := monthsID[t.Month()]
	year := t.Year()

	return fmt.Sprintf("%d %s %d", day, month, year)
}

// TimeAgo menampilkan waktu relatif (Human Readable).
// Contoh: "Baru saja", "5 menit yang lalu".
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	seconds := diff.Seconds()
	minutes := diff.Minutes()
	hours := diff.Hours()
	days := hours / 24

	switch {
	case seconds < 60:
		return "Baru saja"
	case minutes < 60:
		return fmt.Sprintf("%.0f menit yang lalu", minutes)
	case hours < 24:
		return fmt.Sprintf("%.0f jam yang lalu", hours)
	case days < 7:
		return fmt.Sprintf("%.0f hari yang lalu", days)
	case days < 30:
		return fmt.Sprintf("%.0f minggu yang lalu", math.Floor(days/7))
	default:
		return FormatDateID(t) // Jika lebih dari sebulan, tampilkan tanggal lengkap
	}
}
