package main

import (
	"fmt"
	"sync"
)

/*

If Go did not support the sync package, including tools like Mutex or atomic operations,
you would need to rely on other concurrency control mechanisms to avoid the data race.

One effective way to solve this issue would be to use channels to coordinate access
to shared state (i.e., counter). Channels are Go's primary concurrency primitive
and can be used to synchronize operations between goroutines.

*/

func main() {
	counterChan := make(chan int)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	// Counter goroutine that listens for increments
	go func() {
		counter := 0
		for inc := range counterChan {
			counter += inc
		}
		fmt.Println("Final Counter:", counter) // This will print after the channel is closed
		doneChan <- struct{}{}
	}()

	// Launch 1000 goroutines to send increment requests
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counterChan <- 1 // Send increment
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the counter channel to signal completion
	close(counterChan)

	<-doneChan
}
