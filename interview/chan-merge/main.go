package main

import (
	"fmt"
	"sync"
)

// merge N channels into one channel
func mergeChannels(chans ...<-chan int) chan int {
	res := make(chan int)

	go func() {
		var wg sync.WaitGroup
		// wg.Add(len(chans))

		for _, ch := range chans {
			wg.Add(1)

			go func(ch <-chan int, wg *sync.WaitGroup) {
				defer wg.Done()

				for v := range ch {
					res <- v
				}
			}(ch, &wg)
		}

		wg.Wait()
		close(res)
	}()

	return res
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			a <- i
		}
		close(a)
	}()

	go func() {
		for i := 10; i < 20; i++ {
			b <- i
		}
		close(b)
	}()

	go func() {
		for i := 20; i < 30; i++ {
			c <- i
		}
		close(c)
	}()

	res := make([]int, 0)
	for num := range mergeChannels(a, b, c) {
		fmt.Println(num)
		res = append(res, num)
	}

	fmt.Println(len(res) == 30)
}
