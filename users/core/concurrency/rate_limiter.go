package concurrency

import (
	"sync"
	"time"
)

// SlidingWindowRateLimiter menggunakan ring buffer waktu
// untuk menghitung rate limit secara presisi.
type SlidingWindowRateLimiter struct {
	mu     sync.Mutex
	times  []time.Time
	head   int
	count  int
	limit  int
	window time.Duration
}

func NewSlidingWindowRateLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	if limit <= 0 {
		limit = 1
	}
	// allocate limit+1 to differentiate full vs empty buffer
	return &SlidingWindowRateLimiter{
		times:  make([]time.Time, limit+1),
		limit:  limit,
		window: window,
	}
}

func (rl *SlidingWindowRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// purge expired entries (membersihkan data lama)
	for rl.count > 0 {
		t := rl.times[rl.head]
		if t.After(cutoff) {
			break
		}
		// Move head forward (circular buffer)
		rl.head = (rl.head + 1) % len(rl.times)
		rl.count--
	}

	// Cek apakah kuota penuh
	if rl.count >= rl.limit {
		return false
	}

	// Masukkan timestamp request baru
	idx := (rl.head + rl.count) % len(rl.times)
	rl.times[idx] = now
	rl.count++
	return true
}

func (rl *SlidingWindowRateLimiter) GetStats() (int, float64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	pct := 0.0
	if rl.limit > 0 {
		pct = float64(rl.count) / float64(rl.limit) * 100
	}
	// Mengembalikan (Jumlah Request Saat Ini, Persentase Penggunaan)
	return rl.count, pct
}
