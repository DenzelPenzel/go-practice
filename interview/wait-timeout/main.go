/*

Create WaitTimeout function for the Wait function in Sync.WaitGroup

Implement WaitTimeout function:
	- requires adding timeout functionality to sync.WaitGroup
	- if the timeout has been reached, return true
	- if the WaitGroup naturally completes, return false
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	c := make(chan struct{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int, close <-chan struct{}) {
			defer wg.Done()
			<-close
			fmt.Println(num)
		}(i, c)
	}

	if WaitTimeout(&wg, time.Second*5) {
		close(c)
		fmt.Println("timeout exit")
	}

	time.Sleep(time.Second * 10)
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	// sync.WaitGroup the object Wait() func is blocking
	// time.Timer also blocking
	// when two objects are concurrently blocked,
	// each of them can initiates a separate coroutine to handle their respective blockings
	// the challenge lies in determining which of these blockings is resolved first
	// solution - use unbuffered 1 channel
	ch := make(chan bool, 1)

	go time.AfterFunc(timeout, func() {
		// timeout happened, return true
		ch <- true
	})

	go func() {
		wg.Wait()
		// Wait group finished, return false
		ch <- false
	}()

	return <-ch
}
