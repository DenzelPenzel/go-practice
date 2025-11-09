package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

type Event struct {
	ID   int
	Data string
}

type Batcher struct {
	group *errgroup.Group
	gctx  context.Context

	sink Sink
	ins  chan Event

	closing  chan struct{}
	done     chan struct{}
	stopOnce sync.Once
}

func NewBatcher(sink Sink) *Batcher {
	return &Batcher{
		sink:    sink,
		ins:     make(chan Event, 100),
		closing: make(chan struct{}),
		done:    make(chan struct{}),
	}
}

func (b *Batcher) initiateShutdown() {
	b.stopOnce.Do(func() {
		close(b.closing)
		close(b.ins)
	})
}

func (b *Batcher) Start(ctx context.Context) {
	workers := runtime.NumCPU()
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)
	b.group = g
	b.gctx = gctx

	go b.loop(gctx)
}

func (b *Batcher) loop(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	defer close(b.done)
	defer b.initiateShutdown()

	buf := make([]Event, 0, 100)

	flush := func(batch []Event) {
		b.group.Go(func() error {
			fmt.Println("flushing batch", len(batch))
			fctx, fcancel := context.WithTimeout(b.gctx, 10*time.Second)
			defer fcancel()
			return b.sink.write(fctx, batch)
		})
	}

	for {
		select {
		case <-ctx.Done():
			if len(buf) > 0 {
				batch := append([]Event(nil), buf...)
				flush(batch)
			}
			return

		case ev, ok := <-b.ins:
			if !ok {
				if len(buf) > 0 {
					batch := append([]Event(nil), buf...)
					flush(batch)
				}
				return
			}

			buf = append(buf, ev)
			if len(buf) >= 100 {
				batch := append([]Event(nil), buf...)
				buf = buf[:0]
				flush(batch)
			}

		case <-ticker.C:
			if len(buf) > 0 {
				batch := append([]Event(nil), buf...)
				buf = buf[:0]
				flush(batch)
			}
		}
	}
}

func (b *Batcher) Put(ctx context.Context, ev Event) error {
	select {
	case <-b.closing:
		return errors.New("batcher is closed")

	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-b.closing:
		return errors.New("batcher is closed")

	case b.ins <- ev:
		return nil
	}
}

func (b *Batcher) Stop(ctx context.Context) error {
	if b.group == nil {
		return errors.New("batcher not started")
	}

	select {
	case <-b.done:
		return errors.New("batcher is closed")

	default:
	}

	b.initiateShutdown()

	select {
	case <-b.done:
	case <-ctx.Done():
		return ctx.Err()
	}

	return b.group.Wait()
}

type Sink interface {
	write(ctx context.Context, batch []Event) error
}

// MockSink simulates a database with variable latency and optional failures.
type MockSink struct {
	MinLatency time.Duration
	MaxLatency time.Duration
	// FailEvery causes a failure every Nth batch (if > 0).
	FailEvery int
	batches   uint64
}

func (m *MockSink) write(ctx context.Context, batch []Event) error {
	// Simulate latency
	lat := m.MinLatency
	if m.MaxLatency > m.MinLatency {
		lat += time.Duration(rand.Int63n(int64(m.MaxLatency - m.MinLatency)))
	}
	select {
	case <-time.After(lat):
	case <-ctx.Done():
		return ctx.Err()
	}

	n := atomic.AddUint64(&m.batches, 1)
	if m.FailEvery > 0 && int(n)%m.FailEvery == 0 {
		return errors.New("mock sink failure")
	}
	// Pretend batch is atomically persisted
	return nil
}

func main() {
	sink := &MockSink{
		MinLatency: 30 * time.Millisecond,
		MaxLatency: 120 * time.Millisecond,
		FailEvery:  0, // set to, e.g., 5 to see error propagation
	}

	b := NewBatcher(sink)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.Start(ctx)

	prodCtx, prodCancel := context.WithTimeout(ctx, 5*time.Second)
	defer prodCancel()

	go func() {
		for i := 0; i < 100; i++ {
			ev := Event{ID: i, Data: fmt.Sprintf("e-%d", i)}
			if err := b.Put(prodCtx, ev); err != nil {
				fmt.Println("Put stopped:", err)
				return
			}
			time.Sleep(time.Duration(10+rand.Intn(15)) * time.Millisecond)

		}
	}()

	time.Sleep(2 * time.Second)

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	if err := b.Stop(stopCtx); err != nil {
		fmt.Println("stop error:", err)
	} else {
		fmt.Println("Done!")
	}
}
