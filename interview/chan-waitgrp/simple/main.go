package main

import (
	"fmt"
	"time"
)

type ChannelWaitGroup struct {
	// chAdd is used to send increments/decrements to the counter.
	chAdd chan int
	// chWait is used to send a request to wait.
	// The sender provides a channel (chan struct{}) that the monitor will close when ready.
	chWait chan chan struct{}
}

func NewChannelWaitGroup() *ChannelWaitGroup {
	wg := &ChannelWaitGroup{
		chAdd:  make(chan int),
		chWait: make(chan chan struct{}),
	}

	go wg.monitor()

	return wg
}

func (wg *ChannelWaitGroup) monitor() {
	var counter int
	var waiters []chan struct{}

	for {
		select {
		case delta := <-wg.chAdd:
			counter += delta

			if counter < 0 {
				panic("sync: negative WaitGroup counter")
			}

			if counter == 0 {
				for _, ch := range waiters {
					close(ch) // Closing the channel unblocks the waiter
				}
				// Reset the waiters list
				waiters = nil
			}

		case ch := <-wg.chWait:
			if counter == 0 {
				close(ch)
			} else {
				waiters = append(waiters, ch)
			}
		}

	}
}

func (wg *ChannelWaitGroup) Add(delta int) {
	wg.chAdd <- delta
}

func (wg *ChannelWaitGroup) Done() {
	wg.chAdd <- -1
}

func (wg *ChannelWaitGroup) Wait() {
	// Create a unique channel for this specific Wait call
	reply := make(chan struct{})

	// Send the channel to the monitor
	wg.chWait <- reply

	// Block here until the monitor closes the 'reply' channel
	<-reply
}

func main() {
	wg := NewChannelWaitGroup()

	workerCount := 3
	fmt.Printf("Starting %d workers...\n", workerCount)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Second * time.Duration(id)) // Simulate work
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	fmt.Println("Main: Waiting for workers...")
	wg.Wait()
	fmt.Println("Main: All workers completed.")
}
