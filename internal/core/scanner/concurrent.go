package scanner

import (
	"context"
	"runtime"
	"sync"
)

const (
	// DEFAULT_WORKER_COUNT uses all available CPU cores
	DEFAULT_WORKER_COUNT = 0
	
	// MAX_WORKER_COUNT limits maximum workers to prevent resource exhaustion
	MAX_WORKER_COUNT = 50
	
	// DEFAULT_JOB_BUFFER_SIZE provides reasonable buffering for job distribution
	DEFAULT_JOB_BUFFER_SIZE = 100
)

// Job represents a directory or file to be processed
type Job struct {
	Path  string
	Depth int
}

// WorkerPool manages concurrent file processing
type WorkerPool struct {
	workerCount int
	jobChan     chan Job
	resultChan  chan ScanResult
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	processor   JobProcessor
}

// JobProcessor defines the interface for processing jobs
type JobProcessor interface {
	ProcessJob(ctx context.Context, job Job) []ScanResult
}

// NewWorkerPool creates a new worker pool with specified configuration
func NewWorkerPool(workerCount int, bufferSize int, processor JobProcessor) *WorkerPool {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}
	if workerCount > MAX_WORKER_COUNT {
		workerCount = MAX_WORKER_COUNT
	}
	if bufferSize <= 0 {
		bufferSize = DEFAULT_JOB_BUFFER_SIZE
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workerCount: workerCount,
		jobChan:     make(chan Job, bufferSize),
		resultChan:  make(chan ScanResult, bufferSize),
		ctx:         ctx,
		cancel:      cancel,
		processor:   processor,
	}
}

// Start initializes and starts the worker goroutines
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobChan)
	wp.wg.Wait()
	close(wp.resultChan)
	wp.cancel()
}

// SubmitJob adds a job to the processing queue
func (wp *WorkerPool) SubmitJob(job Job) {
	select {
	case wp.jobChan <- job:
	case <-wp.ctx.Done():
		// Pool is shutting down
	}
}

// Results returns the channel for receiving scan results
func (wp *WorkerPool) Results() <-chan ScanResult {
	return wp.resultChan
}

// worker processes jobs from the job channel
func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobChan:
			if !ok {
				// Job channel closed, worker should exit
				return
			}
			
			// Process the job
			results := wp.processor.ProcessJob(wp.ctx, job)
			
			// Send results
			for _, result := range results {
				select {
				case wp.resultChan <- result:
				case <-wp.ctx.Done():
					return
				}
			}
			
		case <-wp.ctx.Done():
			// Context cancelled, worker should exit
			return
		}
	}
}

// WorkerCount returns the number of workers in the pool
func (wp *WorkerPool) WorkerCount() int {
	return wp.workerCount
}