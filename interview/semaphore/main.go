package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Semaphore struct {
	sem chan struct{}
	// Optional: To ensure that Close is called only once
	mu     sync.Mutex
	closed bool
}

func NewSemaphore(capacity int) *Semaphore {
	if capacity <= 0 {
		panic("semaphore capacity must be positive")
	}
	return &Semaphore{
		sem: make(chan struct{}, capacity),
	}
}

// Acquire tries to acquire a semaphore permit.
// It blocks until a permit is available or the context is canceled.
// Returns an error if the context is canceled.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire attempts to acquire a semaphore permit without blocking.
// Returns true if a permit was acquired, false otherwise.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases a semaphore permit.
// It panics if there are more releases than acquires
func (s *Semaphore) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		panic("semaphore is closed")
	}
	select {
	case <-s.sem:
		// Successfully released
	default:
		panic("release called more times than acquire")
	}

}

// Close closes the semaphore, releasing any resources.
// After closing, no further Acquire calls should be made.
func (s *Semaphore) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.closed {
		close(s.sem)
		s.closed = true
	}
}

// Make a custom waitGroup on a semaphore

/*
We want to make a semaphore that will wait for five goroutines to complete!

1. Create buffered channel, inside each goroutine we put a value in it
2. At the end we will expect that everything is ok - we will subtract all the values ​​from the channel
*/
func main() {
	sem := NewSemaphore(3)
	defer sem.Close()

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel() // Ensures that resources are freed when the goroutine completes

			if err := sem.Acquire(ctx); err != nil {
				log.Printf("Goroutine %d: failed to acquire semaphore: %v", id, err)
				return
			}

			// Ensure that the permit is released
			defer sem.Release()

			// Simulate work
			log.Printf("Goroutine %d: acquired semaphore", id)
			time.Sleep(1 * time.Second)
			log.Printf("Goroutine %d: released semaphore", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("Done!")
}
