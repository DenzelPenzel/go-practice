package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var ErrSemaphoreClosed = errors.New("semaphore closed")

type Semaphore struct {
	permits chan struct{}
	closed  chan struct{}
	once    sync.Once
}

func NewSemaphore(size int) *Semaphore {
	if size <= 0 {
		panic("semaphore capacity must be positive")
	}

	permits := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		permits <- struct{}{}
	}

	return &Semaphore{
		permits: permits,
		closed:  make(chan struct{}),
	}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("acquire permit: %w", ctx.Err())

	case <-s.closed:
		return ErrSemaphoreClosed

	case <-s.permits:
		if s.isClosed() {
			s.Release()
			return ErrSemaphoreClosed
		}
		return nil
	}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case <-s.permits:
		if s.isClosed() {
			s.Release()
			return false
		}
		return true

	default:
		return false
	}
}

func (s *Semaphore) Release() {
	select {
	case s.permits <- struct{}{}:
		return

	default:
		panic("semaphore: release without matching acquire")
	}
}

func (s *Semaphore) Close() {
	s.once.Do(func() {
		close(s.closed)
	})
}

func (s *Semaphore) isClosed() bool {
	select {
	case <-s.closed:
		return true

	default:
		return false
	}
}

const (
	maxConcurrentWorkers = 3
	totalTasks           = 10
	acquireTimeout       = 2 * time.Second
	simulatedWork        = time.Second
)

func main() {
	sem := NewSemaphore(maxConcurrentWorkers)
	defer sem.Close()

	var wg sync.WaitGroup
	ctx := context.Background()

	for id := 1; id <= totalTasks; id++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			runWorker(ctx, sem, id)
		}(id)
	}

	wg.Wait()
	fmt.Println("Done!")
}

func runWorker(ctx context.Context, sem *Semaphore, id int) {
	ctx, cancel := context.WithTimeout(ctx, acquireTimeout)
	defer cancel()

	if err := sem.Acquire(ctx); err != nil {
		switch {
		case errors.Is(err, ErrSemaphoreClosed):
			log.Printf("worker %d: semaphore closed", id)

		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("worker %d: acquire timed out", id)

		default:
			log.Printf("worker %d: unable to acquire permit: %v", id, err)
		}
		return
	}
	defer sem.Release()

	log.Printf("worker %d: acquired permit, execute work", id)

	select {
	case <-time.After(simulatedWork):
		log.Printf("worker %d: completed work", id)

	case <-ctx.Done():
		log.Printf("worker %d: context cancelled: %v", id, ctx.Err())
	}
}
