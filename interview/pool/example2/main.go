package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MaxPoolBufSize ensures we don't keep massive buffers in memory.
// If a buffer grows too large (e.g., 100MB) due to a specific edge case,
// we want to drop it so the GC can reclaim that memory, rather than
// holding onto it forever in the pool
const MaxPoolBufSize = 64 * 1024 // 64KB

var bufPool = sync.Pool{
	New: func() interface{} {
		// We pre-allocate a small buffer to avoid immediate resizing.
		// bytes.Buffer is the ideal candidate for pooling because it contains
		// an underlying []byte slice that causes GC pressure if re-allocated constantly
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

// GetBuffer fetches a buffer from the pool
func GetBuffer() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	// CRITICAL STEP 1: Reset the buffer.
	// Reset() retains the underlying storage but sets length to 0.
	// If you forget this, you append new data to old data (data corruption).
	buf.Reset()

	// CRITICAL STEP 2: Prevent Memory Pinning.
	// If the buffer grew huge ( > 64KB), do NOT put it back.
	// Let the Garbage Collector reclaim it. Otherwise, your application's
	// base memory usage will ratchet up to the size of the largest request ever processed.
	if buf.Cap() > MaxPoolBufSize {
		return
	}

	bufPool.Put(buf)
}

type LogEntry struct {
	Timestamp string `json:"ts"`
	Level     string `json:"level"`
	Message   string `json:"msg"`
	RequestID string `json:"req_id"`
}

func LogInfo(reqID, msg string) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   msg,
		RequestID: reqID,
	}

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(entry); err != nil {
		fmt.Println("Error encoding:", err)
		return
	}

	fmt.Print(buf.String())
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("Starting high-concurrency logging...")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				reqID := fmt.Sprintf("req-%d-%d", id, j)
				LogInfo(reqID, "User logged in")
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Done!")
}
