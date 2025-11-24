package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ChannelWaitGroup struct {
	counter int64
	// waitCh acts as the barrier. It is closed when counter reaches 0.
	waitCh chan struct{}
	// lock acts as a Mutex to protect the counter and waitCh updates.
	// We use a buffered channel of size 1.
	lock chan struct{}
}

func NewChannelWaitGroup() *ChannelWaitGroup {
	wg := &ChannelWaitGroup{
		lock: make(chan struct{}, 1),
	}
	wg.lock <- struct{}{}
	return wg
}

func (wg *ChannelWaitGroup) Add(delta int) {
	// Acquire Lock
	<-wg.lock
	defer func() {
		wg.lock <- struct{}{} // Release Lock
	}()

	if wg.counter == 0 && delta > 0 {
		// Transitioning from 0 to 1: We need a new wait channel (generation)
		// If there was an old closed channel, replace it.
		wg.waitCh = make(chan struct{})
	}

	wg.counter += int64(delta)

	if wg.counter < 0 {
		panic("sync: negative WaitGroup counter")
	}

	// Transitioning to 0: Broadcast to all waiters
	if wg.counter == 0 && wg.waitCh != nil {
		close(wg.waitCh)
		wg.waitCh = nil
	}
}

func (wg *ChannelWaitGroup) Done() {
	wg.Add(-1)
}

func (wg *ChannelWaitGroup) Wait() {
	if atomic.LoadInt64(&wg.counter) == 0 {
		return
	}

	// Acquire Lock to safely get the current wait channel
	<-wg.lock
	currentCh := wg.waitCh
	currentCounter := wg.counter
	wg.lock <- struct{}{} // Release Lock immediately

	// If counter dropped to 0 between the atomic check and the lock, return.
	if currentCounter == 0 || currentCh == nil {
		return
	}

	// Block on the channel.
	// When Add(-1) reduces counter to 0, it will close this channel.
	<-currentCh
}

func main() {
	wg := NewChannelWaitGroup()

	fmt.Println("Starting workers...")

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(100*id))
			fmt.Printf("Worker %d finished\n", id)
		}(i)
	}

	fmt.Println("Main: Waiting...")
	wg.Wait()
	fmt.Println("Main: All done.")
}
