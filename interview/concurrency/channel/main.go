/*
In Go, once you close a channel using close(ch), you can still read from the channel, but you need to be aware of the implications:

If you try to read from a closed channel, it won't block, and it will immediately return the zero value
of the channel's type.

For example, if you have a channel of type int, attempting to read from the closed channel will return 0 for an integer.

After closing a channel, you can continue to read any remaining values that were sent to the channel before it was closed.
The channel will only be closed for writes, not for reads.
*/

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	test1()
}

func test1() {
	ch := make(chan int)

	go func() {
		ch <- 1
		ch <- 2
		close(ch)
	}()

	for {
		val, ok := <-ch

		if !ok {
			fmt.Println("Channel is closed. ->", val)
			return
		}
		fmt.Println("Received ->", val)
	}
}

// src/runtime/chan.go

// Write to the chan
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	// if chan is closed, at this time, writing is executed, the source code prompts directly panic
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}
}

// Read from chan
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	lock(&c.lock)

	if c.closed != 0 { // chan is closed
		if c.qcount == 0 {
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			unlock(&c.lock)

			// typedmemclr will clean up the response memory according to the type
			// This explains why the closed chan returns a zero value of the corresponding type
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			// reading closing channel
			// channel return: val, ok := <- chan
			return true, false
		}
		// The channel has been closed, but the channel's buffer have data.
	} else {
		// Just found waiting sender with not closed.
		if sg := c.sendq.dequeue(); sg != nil {
			// Found a waiting sender. If buffer is size 0, receive value
			// directly from sender. Otherwise, receive from head of queue
			// and add sender's value to the tail of the queue (both map to
			// the same buffer slot because the queue is full).
			recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
			return true, true
		}
	}
}

/*
Using a buffered channel in this task helps to prevent goroutines from getting stuck due to potential synchronization issues.

Blocking Goroutines:
	if you use an unbuffered channel, the ch <- elem operation in the goroutine will block
	until there is a corresponding <-ch operation on the receiving side
	if there is no other goroutine currently ready to receive from the channel
	the sending goroutine will be blocked, and the iteration over the set will also be blocked

Buffered Channel:
	using a buffered channel with a capacity equal to the length of the set allows the sending goroutine
	to continue sending elements into the channel even if there is no immediate receiver.
	the goroutine won't block until the channel is full (i.e., its buffer is exhausted).

Avoiding Goroutine Deadlock:
	if the channel is unbuffered, there's a risk of a deadlock scenario
	if the channel is created but there are no other goroutines ready to receive from it
	the sending goroutine will block forever
	with a buffered channel, the goroutine can proceed to send elements, and even if there are no immediate receivers
	the buffer allows the goroutine to complete the iteration and close the channel

Release the lock before closing the channel to prevent deadlock
*/

func (set *threadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{}, len(set.s)) // <- use a buffered channel

	go func() {
		set.RLock()
		defer set.RUnlock()

		for elem := range set.s {
			ch <- elem
		}

		close(ch)
	}()

	return ch
}
