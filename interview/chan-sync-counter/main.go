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
	doneChan := make(chan bool)

	var wg sync.WaitGroup

	go func() {
		counter := 0
		for {
			select {
			case inc := <-counterChan:
				counter += inc

			case <-doneChan:
				fmt.Println("Final Counter:", counter)
				close(doneChan)
				return
			}
		}
	}()

	// Launch 1000 goroutines to send increment requests
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counterChan <- 1
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Signal the counter goroutine to print the final value and exit
	doneChan <- true

	// Wait for the counter goroutine to finish
	<-doneChan
}
