package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ----------------------------- Domain types -----------------------------
// Consume blocks from a mock RPC (or generator), parse logs, index events to an in-memory store.
// Must handle chain reorgs: on reorg signal, cancel in-flight work and resume from the new canonical head.
type Block struct {
	Number     int64
	Hash       string
	ParentHash string
	Timestamp  time.Time
	Logs       []Log
}

type Log struct {
	Index int
	Data  string
}

type Event struct {
	BlockNumber int64
	BlockHash   string
	LogIndex    int
	Kind        string
	Payload     string
	Timestamp   time.Time
}

// ----------------------------- Store ------------------------------------

// Store keeps events keyed by (blockHash, logIndex) and allows reverting
// everything from a given block number onward.
type Store struct {
	mu        sync.RWMutex
	byKey     map[string]Event // hash#index -> event
	byKind    map[string]int64
	latestNum int64
	// Optional index by block number for faster reverts
	byBlock map[int64][]string // blockNumber -> keys stored for that block
}

func NewStore() *Store {
	return &Store{
		byKey:   make(map[string]Event, 1024),
		byKind:  make(map[string]int64, 16),
		byBlock: make(map[int64][]string, 256),
	}
}

func evKey(hash string, idx int) string { return fmt.Sprintf("%s#%d", hash, idx) }

func (s *Store) Upsert(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := evKey(ev.BlockHash, ev.LogIndex)
	if _, exists := s.byKey[k]; exists {
		return
	}
	s.byKey[k] = ev
	s.byKind[ev.Kind]++
	s.byBlock[ev.BlockNumber] = append(s.byBlock[ev.BlockNumber], k)
	if ev.BlockNumber > s.latestNum {
		s.latestNum = ev.BlockNumber
	}
}

func (s *Store) MarkLatestProcessed(blockNum int64) {
	s.mu.Lock()
	if blockNum > s.latestNum {
		s.latestNum = blockNum
	}
	s.mu.Unlock()
}

func (s *Store) ReorgRevert(fromBlockNum int64) {
	// Remove all events with BlockNumber >= fromBlockNum
	s.mu.Lock()
	defer s.mu.Unlock()
	for bn, keys := range s.byBlock {
		if bn >= fromBlockNum {
			for _, k := range keys {
				if ev, ok := s.byKey[k]; ok {
					if s.byKind[ev.Kind] > 0 {
						s.byKind[ev.Kind]--
					}
					delete(s.byKey, k)
				}
			}
			delete(s.byBlock, bn)
		}
	}
	if s.latestNum >= fromBlockNum {
		s.latestNum = fromBlockNum - 1
	}
}

func (s *Store) Stats() (events int, byKind map[string]int64, latest int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events = len(s.byKey)
	byKind = make(map[string]int64, len(s.byKind))
	for k, v := range s.byKind {
		byKind[k] = v
	}
	latest = s.latestNum
	return
}

// ----------------------------- Mock chain source -------------------------

// Reorg indicates we should restart from NewHead (block number, inclusive).
type Reorg struct {
	NewHead int64
}

// hashBlock deterministically derives a pseudo hash from parent hash and block number.
func hashBlock(parentHash string, number int64) string {
	h := sha256.New()
	h.Write([]byte(parentHash))
	var b [8]byte
	for i := 0; i < 8; i++ {
		b[i] = byte(number >> (8 * i))
	}
	h.Write(b[:])
	sum := h.Sum(nil)
	return "0x" + hex.EncodeToString(sum[:8]) // short hash for display
}

// StreamBlocks emits a canonical chain starting at 'start', with proper parent links.
// Occasionally it simulates a reorg by signaling a rewind to a previous block number.
// The caller should cancel, revert, and restart from the provided NewHead.
func StreamBlocks(ctx context.Context, start int64, blocksCh chan<- Block, reorgCh chan<- Reorg) error {
	cur := start
	// Choose a parent of the start: for demo, deterministic
	parentHash := hashBlock("0xgenesis", start-1)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Random reorg trigger
		if cur > start && rand.Intn(15) == 0 {
			back := int64(1 + rand.Intn(3))
			newHead := cur - back
			if newHead < start {
				newHead = start
			}
			select {
			case reorgCh <- Reorg{NewHead: newHead}:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		h := hashBlock(parentHash, cur)
		blk := Block{
			Number:     cur,
			Hash:       h,
			ParentHash: parentHash,
			Timestamp:  time.Now(),
			Logs: []Log{
				{Index: 0, Data: fmt.Sprintf("transfer:%d", cur)},
				{Index: 1, Data: fmt.Sprintf("mint:%d", cur)},
			},
		}

		select {
		case blocksCh <- blk:
			parentHash = h
			cur++
			time.Sleep(25 * time.Millisecond) // simulate block time
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// ----------------------------- Parser ------------------------------------

func parseLog(_ context.Context, blk Block, lg Log) (Event, error) {
	// Simulate some work; in real life, decode ABI/topics
	time.Sleep(4 * time.Millisecond)
	kind := "unknown"
	switch {
	case len(lg.Data) >= 9 && lg.Data[:9] == "transfer:":
		kind = "transfer"
	case len(lg.Data) >= 5 && lg.Data[:5] == "mint:":
		kind = "mint"
	}
	return Event{
		BlockNumber: blk.Number,
		BlockHash:   blk.Hash,
		LogIndex:    lg.Index,
		Kind:        kind,
		Payload:     lg.Data,
		Timestamp:   blk.Timestamp,
	}, nil
}

// ----------------------------- Indexer -----------------------------------

type Indexer struct {
	Workers int
	Buf     int
	Store   *Store
	Logger  func(string)
}

func NewIndexer(workers, buf int, store *Store) *Indexer {
	if workers <= 0 {
		workers = 4
	}
	if buf <= 0 {
		buf = 128
	}
	return &Indexer{
		Workers: workers,
		Buf:     buf,
		Store:   store,
		Logger:  func(string) {},
	}
}

func (ix *Indexer) logf(format string, a ...any) {
	if ix.Logger != nil {
		ix.Logger(fmt.Sprintf(format, a...))
	}
}

// Run starts indexing from 'start'. On a reorg signal, it cancels in-flight work,
// reverts store state from the indicated block number onward, and restarts.
func (ix *Indexer) Run(ctx context.Context, start int64) error {
	head := start

	for ctx.Err() == nil {
		blocksCh := make(chan Block, ix.Buf)
		reorgCh := make(chan Reorg, 1)
		loopCtx, loopCancel := context.WithCancel(ctx)

		// Start block stream
		go func() {
			_ = StreamBlocks(loopCtx, head, blocksCh, reorgCh)
			close(blocksCh)
		}()

		eventsCh := make(chan Event, ix.Buf)

		// Single writer for store
		var writerWG sync.WaitGroup
		writerWG.Add(1)
		go func() {
			defer writerWG.Done()
			for ev := range eventsCh {
				ix.Store.Upsert(ev)
				ix.Store.MarkLatestProcessed(ev.BlockNumber)
			}
		}()

		// Workers
		g, gctx := errgroup.WithContext(loopCtx)
		g.SetLimit(ix.Workers)

		// Dispatch parse tasks
		dispatchDone := make(chan struct{})

		go func() {
			defer close(dispatchDone)
			for blk := range blocksCh {
				// Optional: validate linkage or record parent hash for tests
				for _, lg := range blk.Logs {
					blk, lg := blk, lg
					g.Go(func() error {
						ev, err := parseLog(gctx, blk, lg)
						if err != nil {
							return err
						}
						select {
						case eventsCh <- ev:
							return nil

						case <-gctx.Done():
							return gctx.Err()
						}
					})
				}
			}
		}()

		var ro Reorg
		var gotReorg bool

		select {
		case ro = <-reorgCh:
			gotReorg = true
			ix.logf("reorg: rewind to block %d", ro.NewHead)
			loopCancel()

		case <-ctx.Done():
			loopCancel()

		case <-dispatchDone:
			loopCancel()
		}

		_ = g.Wait()
		close(eventsCh)
		writerWG.Wait()

		if gotReorg {
			// Revert everything from NewHead (inclusive) onward.
			ix.Store.ReorgRevert(ro.NewHead)
			head = ro.NewHead
			continue
		}
		return ctx.Err()
	}
	return ctx.Err()
}

// ----------------------------- Demo main ---------------------------------

func main() {
	rand.Seed(time.Now().UnixNano())

	store := NewStore()
	ix := NewIndexer(8, 256, store)
	ix.Logger = func(s string) { fmt.Println(time.Now().Format("15:04:05.000"), s) }

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := int64(1000)
	startAt := time.Now()
	err := ix.Run(ctx, start)
	elapsed := time.Since(startAt)

	events, kinds, latest := store.Stats()
	fmt.Println("stopped:", err)
	fmt.Println("elapsed:", elapsed.Truncate(10*time.Millisecond))
	fmt.Println("events indexed:", events, "latest block processed:", latest)
	fmt.Println("by kind:", kinds)
}
