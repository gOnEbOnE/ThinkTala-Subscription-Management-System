package concurrency

import (
	"math"
	"sort"
	"sync"
	"time"
)

// EnhancedMetrics menggunakan Ring Buffer untuk performa tinggi
// tanpa memakan memori tak terbatas.
type EnhancedMetrics struct {
	mu sync.RWMutex

	// Response time ring buffer
	responseTimes []time.Duration
	rtIndex       int
	rtCount       int
	rtCapacity    int

	// Queue time tracking
	queueTimes []time.Duration
	qtIndex    int
	qtCount    int

	// Current metrics (Cached result)
	avgResponseTime time.Duration
	avgQueueTime    time.Duration
	p95ResponseTime time.Duration
	p99ResponseTime time.Duration
}

func NewEnhancedMetrics(capacity int) *EnhancedMetrics {
	if capacity <= 0 {
		capacity = 1024
	}
	return &EnhancedMetrics{
		responseTimes: make([]time.Duration, capacity),
		queueTimes:    make([]time.Duration, capacity),
		rtCapacity:    capacity,
	}
}

func (m *EnhancedMetrics) AddResponseTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.responseTimes[m.rtIndex] = duration
	m.rtIndex = (m.rtIndex + 1) % m.rtCapacity
	if m.rtCount < m.rtCapacity {
		m.rtCount++
	}

	m.updateResponseMetrics()
}

func (m *EnhancedMetrics) AddQueueTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queueTimes[m.qtIndex] = duration
	m.qtIndex = (m.qtIndex + 1) % m.rtCapacity
	if m.qtCount < m.rtCapacity {
		m.qtCount++
	}

	m.updateQueueMetrics()
}

func (m *EnhancedMetrics) updateResponseMetrics() {
	if m.rtCount == 0 {
		return
	}

	// Calculate average
	var total time.Duration
	// Copy slice agar sort tidak mengacak urutan ring buffer asli
	times := make([]time.Duration, m.rtCount)
	for i := 0; i < m.rtCount; i++ {
		times[i] = m.responseTimes[i]
		total += times[i]
	}
	m.avgResponseTime = total / time.Duration(m.rtCount)

	// Calculate percentiles using sort (n log n)
	// Kita batasi sample sort hanya jika data cukup, agar tidak boros CPU
	if m.rtCount >= 20 {
		sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
		n := len(times)

		p95Index := int(math.Ceil(0.95*float64(n))) - 1
		p99Index := int(math.Ceil(0.99*float64(n))) - 1

		if p95Index < 0 {
			p95Index = 0
		}
		if p99Index < 0 {
			p99Index = 0
		}
		if p95Index >= n {
			p95Index = n - 1
		}
		if p99Index >= n {
			p99Index = n - 1
		}

		m.p95ResponseTime = times[p95Index]
		m.p99ResponseTime = times[p99Index]
	}
}

func (m *EnhancedMetrics) updateQueueMetrics() {
	if m.qtCount == 0 {
		return
	}

	var total time.Duration
	for i := 0; i < m.qtCount; i++ {
		total += m.queueTimes[i]
	}
	m.avgQueueTime = total / time.Duration(m.qtCount)
}

// GetStats returns map snapshot of current metrics
func (m *EnhancedMetrics) GetStats() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]any{
		"avg_response_time": m.avgResponseTime.String(),
		"avg_queue_time":    m.avgQueueTime.String(),
		"p95_response_time": m.p95ResponseTime.String(),
		"p99_response_time": m.p99ResponseTime.String(),
		"sample_count":      m.rtCount,
	}
}
