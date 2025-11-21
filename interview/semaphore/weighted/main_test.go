package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphore_FIFO(t *testing.T) {
	// Limit 1.
	sem := NewSemaphore(1)

	// Acquire the only ticket.
	if err := sem.Acquire(context.Background(), 1); err != nil {
		t.Fatalf("Failed to acquire initial ticket: %v", err)
	}

	var ordering []int
	var mu sync.Mutex

	// Start 3 goroutines. They should acquire in order 1, 2, 3.
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 1; i <= 3; i++ {
		id := i
		go func() {
			defer wg.Done()
			// Request 1 ticket.
			err := sem.Acquire(context.Background(), 1)
			if err != nil {
				t.Errorf("Worker %d failed to acquire: %v", id, err)
				return
			}

			mu.Lock()
			ordering = append(ordering, id)
			mu.Unlock()

			sem.Release(1)
		}()
		// Sleep slightly to ensure they hit the lock in order (probabilistically, but reliable enough for simple test)
		time.Sleep(10 * time.Millisecond)
	}

	// Release the initial ticket to start the chain.
	sem.Release(1)

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(ordering) != 3 {
		t.Fatalf("Expected 3 workers, got %d", len(ordering))
	}
	if ordering[0] != 1 || ordering[1] != 2 || ordering[2] != 3 {
		t.Errorf("Expected ordering [1 2 3], got %v", ordering)
	}
}

func TestSemaphore_IllegalRelease(t *testing.T) {
	sem := NewSemaphore(1)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on illegal release, but got none")
		} else {
			// Check if it's the expected error message or error
			// The code currently uses a string variable or string constant
			// ErrIllegalRelease = "..."
			if r != ErrIllegalRelease {
				t.Errorf("Expected ErrIllegalRelease panic, got %v", r)
			}
		}
	}()

	sem.Release(1) // Should panic
}
