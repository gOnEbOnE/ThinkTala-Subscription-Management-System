package concurrency

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool manages a set of goroutines
type WorkerPool struct {
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	queue       chan Job
	workerCount int
	name        string

	// Reference ke "Otak" (Dispatcher) untuk akses Handler & Metrics
	dispatcher *Dispatcher

	// Metrics lokal
	activeJobs    int64
	processedJobs int64
}

// newWorkerPool creates a new pool instance
func newWorkerPool(ctx context.Context, name string, workerCount, queueSize int, d *Dispatcher) *WorkerPool {
	poolCtx, cancel := context.WithCancel(ctx)
	if queueSize <= 0 {
		queueSize = 1
	}

	pool := &WorkerPool{
		ctx:         poolCtx,
		cancel:      cancel,
		queue:       make(chan Job, queueSize),
		workerCount: workerCount,
		name:        name,
		dispatcher:  d,
	}

	for i := 0; i < workerCount; i++ {
		pool.wg.Add(1)
		go pool.workerLoop(fmt.Sprintf("%s-%d", name, i))
	}

	log.Printf("[Concurrency] Pool '%s' started: %d workers, queue: %d", name, workerCount, queueSize)
	return pool
}

func (wp *WorkerPool) workerLoop(id string) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.queue:
			if !ok {
				return
			}
			wp.processJob(id, job)
		}
	}
}

func (wp *WorkerPool) processJob(workerID string, job Job) {
	atomic.AddInt64(&wp.activeJobs, 1)
	defer atomic.AddInt64(&wp.activeJobs, -1)

	// Hitung Queue Time (Metric)
	queueDuration := time.Since(job.Created)
	if wp.dispatcher.metrics != nil {
		wp.dispatcher.metrics.AddQueueTime(queueDuration)
	}

	// Setup Timeout
	timeout := job.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(wp.ctx, timeout)
	defer cancel()

	// Eksekusi Handler
	start := time.Now()
	var resultData any
	var err error
	var panicErr any

	// Safety: Recover from panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
				log.Printf("[Panic] Worker %s on job %s: %v", workerID, job.ID, r)
			}
		}()

		// 1. Cari Handler berdasarkan JobType
		// NOTE: handlers ada di dispatcher.go (package concurrency), jadi bisa diakses langsung
		handler, exists := wp.dispatcher.handlers[job.Type]
		if !exists {
			err = fmt.Errorf("handler not found for job type: %s", job.Type)
			return
		}

		// 2. Jalankan Handler (Logic Bisnis Aplikasi)
		resultData, err = handler(ctx, job.Payload)
	}()

	execDuration := time.Since(start)

	if wp.dispatcher.metrics != nil {
		wp.dispatcher.metrics.AddResponseTime(execDuration)
	}

	// Update Circuit Breaker & Metrics Global
	if wp.dispatcher.circuitBreaker != nil {
		if err != nil || panicErr != nil {
			// Untuk sekarang kita anggap semua panic adalah system failure
			if panicErr != nil {
				wp.dispatcher.circuitBreaker.RecordFailure()
			}
			// Error logic bisnis biasa tidak selalu record failure circuit breaker
		} else {
			wp.dispatcher.circuitBreaker.RecordSuccess()
		}
	}

	atomic.AddInt64(&wp.processedJobs, 1)

	// Kirim hasil balik (Jika ada yang menunggu/sync)
	if job.ResultChan != nil {
		res := JobResult{
			Success: err == nil && panicErr == nil,
			Data:    resultData,
			Error:   err,
		}

		if panicErr != nil {
			res.Error = fmt.Errorf("internal server error (panic): %v", panicErr)
			res.Success = false
		}

		// Non-blocking send (buffer 1)
		select {
		case job.ResultChan <- res:
		default:
		}
		close(job.ResultChan)
	}
}

// Submit job ke antrean pool
func (wp *WorkerPool) submit(job Job) error {
	select {
	case wp.queue <- job:
		return nil
	default:
		return fmt.Errorf("queue full (%s)", wp.name)
	}
}

// Shutdown pool
func (wp *WorkerPool) shutdown() {
	wp.cancel()
	wp.wg.Wait()
}
