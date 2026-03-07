package concurrency

import (
	"context"
	"time"
)

// =================================================================
// TYPE DEFINITIONS
// =================================================================

const (
	PriorityHigh   = 0
	PriorityNormal = 1
	PriorityLow    = 2
)

// JobType alias string agar kompatibel dengan main.go
type JobType string

// HandlerFunc signature untuk worker logic
type HandlerFunc func(ctx context.Context, payload any) (any, error)

// TAMBAHAN PENTING: Alias agar core/app.go mengenali tipe ini
type JobHandler = HandlerFunc

// Job alias untuk worker pool (Internal)
type Job struct {
	ID         string
	Type       JobType
	Payload    any
	Priority   int
	Timeout    time.Duration
	Created    time.Time
	ResultChan chan JobResult
}

// JobResult hasil eksekusi worker
type JobResult struct {
	Data    any
	Error   error
	Success bool
}

// Config untuk Worker Pool
type Config struct {
	HighPriorityWorkers     int
	NormalPriorityWorkers   int
	LowPriorityWorkers      int
	QueueSizePerPriority    int
	MaxConcurrentJobs       int
	QueueFullThreshold      float64
	MaxRetries              int
	JobTimeoutDefault       time.Duration
	ShutdownTimeout         time.Duration
	EnableAdaptiveRateLimit bool
	HealthCheckRate         time.Duration
	MetricsRate             time.Duration
}
