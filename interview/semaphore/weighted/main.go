package main

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrNoTickets is returned by TryAcquire when resources are unavailable.
	ErrNoTickets = errors.New("semaphore: could not acquire requested tickets")
	// ErrIllegalRelease panics if Release is called more times than Acquire.
	ErrIllegalRelease = errors.New("semaphore: release without an acquisition")
)

type waiter struct {
	n     int64         // How many tickets are requested
	ready chan struct{} // Channel to signal when tickets are acquired
}

// Semaphore provides a weighted semaphore implementation
type Semaphore struct {
	size    int64
	cur     int64
	mu      sync.Mutex
	waiters *list.List
}

func NewSemaphore(limit int64) *Semaphore {
	return &Semaphore{
		size:    limit,
		waiters: list.New(),
	}
}

// Acquire blocks until n tickets are available or the context is done.
// On success, returns nil. On failure, returns ctx.Err().
func (s *Semaphore) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()

	// Fast path: If enough tickets are available and no one is waiting, take them.
	if s.waiters.Len() == 0 && s.cur+n <= s.size {
		s.cur += n
		s.mu.Unlock()
		return nil
	}

	// Slow path: Check if context is already done before queuing
	if err := ctx.Err(); err != nil {
		s.mu.Unlock()
		return err
	}

	// Create a waiter and add to the back of the queue (FIFO)
	w := &waiter{
		n:     n,
		ready: make(chan struct{}),
	}
	elem := s.waiters.PushBack(w)
	s.mu.Unlock()

	// Block until ready or context cancelled
	select {
	case <-ctx.Done():
		// Context cancelled. We must remove ourselves from the queue.
		s.mu.Lock()
		defer s.mu.Unlock()

		// If the ready channel was closed *before* we grabbed the lock,
		// it means we actually succeeded in acquiring the semaphore despite the context.
		// We must honor that success to avoid leaking tickets.
		select {
		case <-w.ready:
			return nil
		default:
		}

		s.waiters.Remove(elem)
		// Since we left the queue, it's possible we were blocking others who can now run.
		// Trigger a notification scan.
		s.notifyWaiters()
		return ctx.Err()

	case <-w.ready:
		// Acquired successfully, go here when close the chan w.ready
		return nil
	}
}

// TryAcquire attempts to acquire n tickets without blocking.
// Returns true on success, false if not enough tickets are available.
func (s *Semaphore) TryAcquire(n int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.waiters.Len() == 0 && s.cur+n <= s.size {
		s.cur += n
		return true
	}
	return false
}

// Release releases n tickets back to the semaphore.
func (s *Semaphore) Release(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cur -= n
	if s.cur < 0 {
		panic(ErrIllegalRelease)
	}

	s.notifyWaiters()
}

// notifyWaiters scans the list of waiters and wakes up as many as possible.
// This must be called while holding the lock.
func (s *Semaphore) notifyWaiters() {
	for {
		next := s.waiters.Front()
		if next == nil {
			break // No waiters
		}

		w := next.Value.(*waiter)
		if s.cur+w.n > s.size {
			// Not enough tickets for the next waiter.
			// Since we enforce FIFO, we stop here to prevent starvation
			// of large requests by small requests.
			break
		}

		// Allocate tickets
		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready) // Signal the waiter
	}
}

func main() {
	// Create a semaphore with a limit of 5 connections
	sem := NewSemaphore(5)
	var wg sync.WaitGroup

	// Simulate 10 workers trying to access resources
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Create a context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			// Attempt to acquire 2 tickets (heavy worker)
			weight := int64(2)
			fmt.Printf("Worker %d requesting %d tickets...\n", id, weight)

			if err := sem.Acquire(ctx, weight); err != nil {
				fmt.Printf("Worker %d timed out or cancelled: %v\n", id, err)
				return
			}

			// Critical Section
			fmt.Printf("Worker %d acquired tickets. Working...\n", id)
			time.Sleep(500 * time.Millisecond)

			// Release
			sem.Release(weight)
			fmt.Printf("Worker %d released tickets.\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers done.")
}
