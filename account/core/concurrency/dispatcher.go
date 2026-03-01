package concurrency

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// =================================================================
// DISPATCHER IMPLEMENTATION
// =================================================================

type Dispatcher struct {
	// Pools (Antrian Worker ada di dalam sini)
	highPool   *WorkerPool
	normalPool *WorkerPool
	lowPool    *WorkerPool

	// Registry Handler
	handlers map[JobType]HandlerFunc

	// Components (Advanced Features dari file utils_*.go)
	// Pastikan file utils_cb.go, utils_rl.go, utils_metrics.go ada
	circuitBreaker *EnhancedCircuitBreaker
	rateLimiter    *SlidingWindowRateLimiter
	metrics        *EnhancedMetrics

	// Config & State
	config     Config
	jobCounter int64
}

// NewDispatcher membuat engine baru
func NewDispatcher(cfg Config) *Dispatcher {
	d := &Dispatcher{
		handlers: make(map[JobType]HandlerFunc),
		// Inisialisasi komponen pendukung
		circuitBreaker: NewCircuitBreaker(100, 10*time.Second),
		rateLimiter:    NewSlidingWindowRateLimiter(10000, time.Second),
		metrics:        NewEnhancedMetrics(10000),
		config:         cfg,
	}

	ctx := context.Background() // Context root

	// Init Pools (Mengirim pointer dispatcher 'd' agar worker bisa akses handler)
	d.highPool = newWorkerPool(ctx, "high", cfg.HighPriorityWorkers, cfg.QueueSizePerPriority, d)
	d.normalPool = newWorkerPool(ctx, "normal", cfg.NormalPriorityWorkers, cfg.QueueSizePerPriority, d)
	d.lowPool = newWorkerPool(ctx, "low", cfg.LowPriorityWorkers, cfg.QueueSizePerPriority, d)

	return d
}

// =================================================================
// COMPATIBILITY METHODS (Untuk dipanggil dari app.go)
// =================================================================

// Start: Menyalakan monitoring
func (d *Dispatcher) Start() {
	go d.monitorLoop()
	log.Println("[Concurrency] Dispatcher Started (Monitoring Active)")
}

// Stop: Alias untuk Shutdown
func (d *Dispatcher) Stop() {
	d.Shutdown()
}

// RegisterJobHandler: Wrapper untuk Register yang menerima string
func (d *Dispatcher) RegisterJobHandler(name string, handler func(context.Context, any) (any, error)) {
	d.Register(JobType(name), HandlerFunc(handler))
}

// =================================================================
// CORE METHODS (Logic Utama)
// =================================================================

// Register mendaftarkan fungsi logic bisnis
func (d *Dispatcher) Register(name JobType, handler HandlerFunc) {
	d.handlers[name] = handler
	log.Printf("[Concurrency] Handler registered: %s", name)
}

// Dispatch: Fire and Forget (Kirim dan lupakan)
func (d *Dispatcher) Dispatch(ctx context.Context, jobType JobType, payload any, priority int) (string, error) {
	return d.enqueue(ctx, jobType, payload, priority, nil)
}

// DispatchAndWait: Kirim dan Tunggu hasil (Blocking)
// Menerima string agar mudah dipanggil dari Controller
func (d *Dispatcher) DispatchAndWait(ctx context.Context, jobTypeStr string, payload any, priority int) (any, error) {
	jobType := JobType(jobTypeStr)

	resultChan := make(chan JobResult, 1)

	_, err := d.enqueue(ctx, jobType, payload, priority, resultChan)
	if err != nil {
		return nil, err
	}

	// Tunggu hasil
	select {
	case res := <-resultChan:
		if res.Error != nil {
			return nil, res.Error
		}
		return res.Data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Internal enqueue logic
func (d *Dispatcher) enqueue(ctx context.Context, jobType JobType, payload any, priority int, resChan chan JobResult) (string, error) {
	// 1. Global Rate Limit Check
	if !d.rateLimiter.Allow() {
		return "", errors.New("rate limit exceeded")
	}

	// 2. Circuit Breaker Check
	if !d.circuitBreaker.Allow() {
		return "", errors.New("service temporarily unavailable (circuit open)")
	}

	// 3. Handler Check
	if _, ok := d.handlers[jobType]; !ok {
		return "", fmt.Errorf("unknown job type: %s", jobType)
	}

	// 4. Create Job
	jobID := fmt.Sprintf("job-%d", atomic.AddInt64(&d.jobCounter, 1))
	job := Job{
		ID:         jobID,
		Type:       jobType,
		Payload:    payload,
		Priority:   priority,
		Timeout:    d.config.JobTimeoutDefault,
		Created:    time.Now(),
		ResultChan: resChan,
	}

	// 5. Select Pool
	var pool *WorkerPool
	switch priority {
	case PriorityHigh:
		pool = d.highPool
	case PriorityLow:
		pool = d.lowPool
	default:
		pool = d.normalPool
	}

	// 6. Submit ke Pool
	if err := pool.submit(job); err != nil {
		return "", err
	}

	return jobID, nil
}

// Shutdown gracefully
func (d *Dispatcher) Shutdown() {
	log.Println("[Concurrency] Shutting down dispatcher...")

	timeout := d.config.ShutdownTimeout
	log.Printf("[Concurrency] Waiting for workers to finish (timeout: %s)...", timeout)

	d.highPool.shutdown()
	d.normalPool.shutdown()
	d.lowPool.shutdown()

	log.Println("[Concurrency] Dispatcher stopped")
}

func (d *Dispatcher) monitorLoop() {
	ticker := time.NewTicker(d.config.HealthCheckRate)
	defer ticker.Stop()
	for range ticker.C {
		d.circuitBreaker.UpdateErrorRate()
	}
}
