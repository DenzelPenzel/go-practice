package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var max int64
	wg := sync.WaitGroup{}

	for i := 1000; i > 0; i-- {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			if i%2 == 0 {
				currentMax := atomic.LoadInt64(&max)
				if i > currentMax {
					atomic.CompareAndSwapInt64(&max, currentMax, i)
				}
			}
		}(int64(i))
	}

	wg.Wait()

	fmt.Printf("Maximum is %d", max)
}
