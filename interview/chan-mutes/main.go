package main

import (
	"fmt"
	"sync"
)

// implement a sync.Mutex using channels

type ChannelMutex struct {
	ch chan struct{}
}

func NewChannelMutex() *ChannelMutex {
	return &ChannelMutex{
		ch: make(chan struct{}, 1), // buffered channel to avoid deadlock
	}
}

// Lock blocks until the mutex is available
func (cm *ChannelMutex) Lock() {
	cm.ch <- struct{}{}
}

// Unlock releases the mutex
func (cm *ChannelMutex) Unlock() {
	<-cm.ch
}

func main() {
	const incrementsPerGoroutine = 100

	var counter int
	var wg sync.WaitGroup
	mu := NewChannelMutex()

	numGoroutines := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("want: %d\n", numGoroutines*incrementsPerGoroutine)
	fmt.Printf("got: %d\n", counter)
}
