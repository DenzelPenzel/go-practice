/*
In Go, once you close a channel using close(ch), you can still read from the channel, but you need to be aware of the implications:

If you try to read from a closed channel, it won't block, and it will immediately return the zero value
of the channel's type.

For example, if you have a channel of type int, attempting to read from the closed channel will return 0 for an integer.

After closing a channel, you can continue to read any remaining values that were sent to the channel before it was closed.
The channel will only be closed for writes, not for reads.
*/

package main

import "fmt"

func main() {
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
