package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

type captureSink struct {
	mu      sync.Mutex
	batches [][]Event
}

func newCaptureSink() *captureSink {
	return &captureSink{}
}

func (c *captureSink) write(ctx context.Context, batch []Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cp := append([]Event(nil), batch...)
	c.batches = append(c.batches, cp)
	return nil
}

func (c *captureSink) batchesSnapshot() [][]Event {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([][]Event, len(c.batches))
	for i, batch := range c.batches {
		out[i] = append([]Event(nil), batch...)
	}
	return out
}

func TestBatcherFlushesOnCapacity(t *testing.T) {
	sink := newCaptureSink()

	b := NewBatcher(sink)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.Start(runCtx)

	for i := 0; i < 100; i++ {
		if err := b.Put(context.Background(), Event{ID: i, Data: "payload"}); err != nil {
			t.Fatalf("put %d: %v", i, err)
		}
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
	defer stopCancel()

	if err := b.Stop(stopCtx); err != nil {
		t.Fatalf("stop: %v", err)
	}

	batches := sink.batchesSnapshot()
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}

	if len(batches[0]) != 100 {
		t.Fatalf("expected batch size 100, got %d", len(batches[0]))
	}

	for i, ev := range batches[0] {
		if ev.ID != i {
			t.Fatalf("batch[%d] has ID %d, want %d", i, ev.ID, i)
		}
	}
}

func TestBatcherFlushesOnStop(t *testing.T) {
	sink := newCaptureSink()
	b := NewBatcher(sink)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.Start(runCtx)

	for i := 0; i < 3; i++ {
		if err := b.Put(context.Background(), Event{ID: i}); err != nil {
			t.Fatalf("put %d: %v", i, err)
		}
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
	defer stopCancel()

	if err := b.Stop(stopCtx); err != nil {
		t.Fatalf("stop: %v", err)
	}

	batches := sink.batchesSnapshot()
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}

	if len(batches[0]) != 3 {
		t.Fatalf("expected batch size 3, got %d", len(batches[0]))
	}
}

func TestBatcherRejectsAfterStop(t *testing.T) {
	sink := newCaptureSink()
	b := NewBatcher(sink)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.Start(runCtx)

	if err := b.Put(context.Background(), Event{ID: 1}); err != nil {
		t.Fatalf("put: %v", err)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
	defer stopCancel()

	if err := b.Stop(stopCtx); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if err := b.Put(context.Background(), Event{ID: 2}); err == nil {
		t.Fatal("expected error on Put after Stop, got nil")
	}
}
