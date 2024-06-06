/*
Suppose there is an extremely long slice, the element type of the slice is int,
and the elements in the slice are arranged in random order.

Time limit is 5 seconds, use multiple goroutines to find whether a
given value exists in the slice, and end the execution of all goroutines immediately
after finding the target value or timing out.

For example, the slice is: [23, 32, 78, 43, 76, 65, 345, 762, ... 915, 86],
and the target value found is 345.

If the target value exists in the slice, the program outputs: "Found it!"
and Immediately cancel the search task that is still executing goroutine.

If the target value is not found within the timeout period, the program outputs:
"Timeout! Not Found" and immediately cancels the search task that is still being executed goroutine.

The most common and often mentioned design pattern in the Go language -
do not communicate through shared memory, but share memory through communication.

https://mp.weixin.qq.com/s/GhC2WDw3VHP91DrrFVCnag

*/

package main

import (
	"context"
	"fmt"
	"time"
)

var data = []int{1, 2, 3, 10, 999, 8, 345, 7, 98, 33, 66, 77, 88, 68, 96}

func main() {
	timer := time.NewTimer(time.Second * 5)
	ctx, cancel := context.WithCancel(context.Background())
	resChan := make(chan bool)
	batch := 3
	target := 345

	for i := 0; i < len(data); i += batch {
		end := i + batch
		if end >= len(data) {
			end = len(data) - 1
		}
		// run G
		go search(ctx, data[i:end], target, resChan)
	}

	select {
	case <-timer.C:
		fmt.Println("Timeout! Not Found")
		cancel()
	case <-resChan:
		fmt.Println("Found it!")
		cancel()
	}

	time.Sleep(time.Second * 2)
}

func search(ctx context.Context, data []int, target int, resChan chan bool) {
	for _, v := range data {
		select {
		case <-ctx.Done():
			fmt.Println("Task canceled!")
			return
		default:
		}
		fmt.Printf("val: %d \n", v)
		time.Sleep(time.Microsecond * 1500)
		if target == v {
			resChan <- true
			return
		}
	}

}
