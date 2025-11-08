package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphoreAcquireRelease(t *testing.T) {
	t.Parallel()

	const capacity = 2
	sem := NewSemaphore(capacity)

	ctx := context.Background()
	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("second acquire failed: %v", err)
	}

	blocked := make(chan struct{})
	go func() {
		defer close(blocked)
		if err := sem.Acquire(ctx); err != nil {
			t.Errorf("unexpected acquire failure: %v", err)
		}
	}()

	select {
	case <-blocked:
		t.Fatal("third acquire should block until a permit is released")
	case <-time.After(50 * time.Millisecond):
	}

	sem.Release()

	select {
	case <-blocked:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("unblocked acquire should succeed after release")
	}
}

func TestSemaphoreTryAcquire(t *testing.T) {
	t.Parallel()

	const capacity = 1
	sem := NewSemaphore(capacity)

	if ok := sem.TryAcquire(); !ok {
		t.Fatal("expected first TryAcquire to succeed")
	}

	if ok := sem.TryAcquire(); ok {
		t.Fatal("expected second TryAcquire to fail when no permits available")
	}

	sem.Release()

	if ok := sem.TryAcquire(); !ok {
		t.Fatal("expected TryAcquire to succeed after release")
	}
}

func TestSemaphoreAcquireContextCancellation(t *testing.T) {
	t.Parallel()

	const capacity = 1
	sem := NewSemaphore(capacity)

	if err := sem.Acquire(context.Background()); err != nil {
		t.Fatalf("pre-acquire failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	if err := sem.Acquire(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Fatalf("acquire returned too quickly: %v", elapsed)
	}
}

func TestSemaphoreClose(t *testing.T) {
	t.Parallel()

	const capacity = 1
	sem := NewSemaphore(capacity)

	if err := sem.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}

	sem.Close()
	sem.Close() // closing twice must not panic

	sem.Release()

	if err := sem.Acquire(context.Background()); !errors.Is(err, ErrSemaphoreClosed) {
		t.Fatalf("expected ErrSemaphoreClosed, got %v", err)
	}

	if ok := sem.TryAcquire(); ok {
		t.Fatal("expected TryAcquire to fail on closed semaphore")
	}
}

func TestSemaphoreReleaseWithoutAcquirePanics(t *testing.T) {
	t.Parallel()

	sem := NewSemaphore(1)

	func() {
		defer func() {
			fmt.Println("recovering...")
			if r := recover(); r == nil {
				t.Fatal("expected panic on extra release")
			}
		}()
		sem.Release()
	}()
}

func TestSemaphoreParallelismLimit(t *testing.T) {
	t.Parallel()

	const (
		capacity   = 3
		iterations = 20
	)
	sem := NewSemaphore(capacity)

	var running atomic.Int32
	var maxSeen atomic.Int32

	errCh := make(chan error, iterations)
	doneCh := make(chan struct{})

	go func() {
		for i := 0; i < iterations; i++ {
			sem.Release()
		}
	}()

	for i := 0; i < iterations; i++ {
		go func() {
			if err := sem.Acquire(context.Background()); err != nil {
				errCh <- err
				return
			}
			defer sem.Release()

			v := running.Add(1)
			for {
				old := maxSeen.Load()
				if v <= old {
					break
				}
				if maxSeen.CompareAndSwap(old, v) {
					break
				}
			}

			time.Sleep(5 * time.Millisecond)
			running.Add(-1)
			errCh <- nil
		}()
	}

	go func() {
		for i := 0; i < iterations; i++ {
			if err := <-errCh; err != nil {
				t.Errorf("unexpected acquire error: %v", err)
			}
		}
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case <-time.After(2 * time.Second):
		t.Fatal("test timed out")
	}

	if maxSeen.Load() > capacity {
		t.Fatalf("observed running workers %d exceeding capacity %d", maxSeen.Load(), capacity)
	}
}
