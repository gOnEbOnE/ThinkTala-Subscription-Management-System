package concurrency

import (
	"sync"
	"sync/atomic"
	"time"
)

type EnhancedCircuitBreaker struct {
	mu              sync.RWMutex
	failures        int64
	successes       int64
	threshold       int64
	recoveryTime    time.Duration
	lastFailTime    time.Time
	state           int32 // 0=closed, 1=open, 2=half-open
	halfOpenAllowed int32

	windowStart    time.Time
	windowDuration time.Duration
	errorRate      float64
}

func NewCircuitBreaker(threshold int64, recoveryTime time.Duration) *EnhancedCircuitBreaker {
	return &EnhancedCircuitBreaker{
		threshold:      threshold,
		recoveryTime:   recoveryTime,
		windowDuration: 60 * time.Second,
		windowStart:    time.Now(),
	}
}

func (cb *EnhancedCircuitBreaker) Allow() bool {
	state := atomic.LoadInt32(&cb.state)
	now := time.Now()

	switch state {
	case 0: // Closed
		return true
	case 1: // Open
		cb.mu.RLock()
		shouldRecover := now.Sub(cb.lastFailTime) > cb.recoveryTime
		cb.mu.RUnlock()

		if shouldRecover && atomic.CompareAndSwapInt32(&cb.state, 1, 2) {
			atomic.StoreInt32(&cb.halfOpenAllowed, 1)
			return true
		}
		return false
	case 2: // Half-open
		return atomic.CompareAndSwapInt32(&cb.halfOpenAllowed, 1, 0)
	default:
		return false
	}
}

func (cb *EnhancedCircuitBreaker) RecordSuccess() {
	atomic.AddInt64(&cb.successes, 1)
	atomic.StoreInt64(&cb.failures, 0)
	atomic.StoreInt32(&cb.state, 0)
	atomic.StoreInt32(&cb.halfOpenAllowed, 0)
}

func (cb *EnhancedCircuitBreaker) RecordFailure() {
	failures := atomic.AddInt64(&cb.failures, 1)
	if failures >= cb.threshold {
		cb.mu.Lock()
		cb.lastFailTime = time.Now()
		cb.mu.Unlock()
		atomic.StoreInt32(&cb.state, 1) // Open
	}
	atomic.StoreInt32(&cb.halfOpenAllowed, 0)
}

func (cb *EnhancedCircuitBreaker) UpdateErrorRate() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	now := time.Now()
	if now.Sub(cb.windowStart) >= cb.windowDuration {
		total := atomic.LoadInt64(&cb.successes) + atomic.LoadInt64(&cb.failures)
		if total > 0 {
			cb.errorRate = float64(atomic.LoadInt64(&cb.failures)) / float64(total) * 100
		}
		atomic.StoreInt64(&cb.successes, 0)
		atomic.StoreInt64(&cb.failures, 0)
		cb.windowStart = now
	}
}
