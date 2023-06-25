package async

import "sync"

type WorkerPoolConfig struct {
	WorkersCount int
	QueueSize    int
}

func WorkerPoolDefaultConfig() WorkerPoolConfig {
	return WorkerPoolConfig{
		WorkersCount: 5,
		QueueSize:    10,
	}
}

type WorkerPool struct {
	// Task queue
	queue chan func()
	// Wait group to ensure all workers finished their tasks after pool is closed
	wg sync.WaitGroup
	// Flag to stop receiving new tasks to queue
	closed bool
	// Mutex granting thread safe access to closed flag
	mu sync.RWMutex
	// Object to perform workerpool closure once
	stopOnce sync.Once
}

func NewWorkerPool(cfg WorkerPoolConfig) *WorkerPool {
	pool := WorkerPool{
		queue: make(chan func(), cfg.QueueSize),
	}

	pool.wg.Add(cfg.WorkersCount)
	for w := 0; w < cfg.WorkersCount; w++ {
		go func() {
			defer pool.wg.Done()
			pool.startWorker()
		}()
	}

	return &pool
}

func (p *WorkerPool) Enqueue(task func()) {
	// Check if pool is still open and able to enqueue task
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return
	}
	p.mu.RUnlock()

	p.queue <- task
}

func (p *WorkerPool) Close() {
	p.stopOnce.Do(func() {
		// Close channel with tasks
		p.mu.Lock()
		p.closed = true
		close(p.queue)
		p.mu.Unlock()
		// Wait for all workers to finish executing
		p.wg.Wait()
	})
}

func (p *WorkerPool) startWorker() {
	for task := range p.queue {
		task()
	}
}
