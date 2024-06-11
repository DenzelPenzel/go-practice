package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Pool

// Issue here that in single-core CPU, the memory may be stable 256MB, but in a multi-core CPU, it may skyrocket

func test() {
	// Pool - mechanism for managing and reusing resources, such as goroutines or objects, to improve performance
	// Why to use?
	// - pool of goroutines avoid the overhead of creating a new goroutine for each task
	// - pool of objects (like database connections or other resources) are reused rather than being created and destroyed frequently
	var pool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	processRequest := func(size int) {
		b := pool.Get().(*bytes.Buffer)
		time.Sleep(500 * time.Millisecond)
		b.Grow(size)
		pool.Put(b)
		time.Sleep(1 * time.Microsecond)
	}

	go func() {
		for {
			processRequest(1 << 28) //  256MiB
		}
	}()

	for i := 0; i < 100; i++ {
		go func() {
			for {
				processRequest(1 << 10) // 1KiB
			}
		}()
	}

	var stats runtime.MemStats
	for i := 0; ; i++ {
		runtime.ReadMemStats(&stats)
		fmt.Printf("Cycle %d: %dB\n", i, stats.Alloc)
		time.Sleep(time.Second)
		runtime.GC()
	}
}

type Data struct {
	A int
	B string
}

var dataPool = sync.Pool{
	New: func() any {
		return &Data{}
	},
}

func main() {
	data := dataPool.Get().(*Data)
	data.A = 10
	data.B = "hello world"

	fmt.Println(data)

	dataPool.Put(data) // Reuse data

	s := make([]int, 0, 100) // Preallocate slice
	for i := 0; i < 100; i++ {
		s = append(s, i)
	}

	fmt.Println(s)
}
