/*
Write code to implement two goroutines, one of which generates
random numbers and writes them to the go channel, and the other
reads numbers from the channel and prints them to the standard output.

The final output is five random numbers.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func sol() {
	out := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer func() {
			close(out)
			wg.Done()
		}()
		for i := 0; i < 5; i++ {
			out <- rand.Intn(5)
		}
	}()

	go func() {
		defer wg.Done()
		for v := range out {
			fmt.Println(v)
		}
	}()

	wg.Wait()
}

func sol_II() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for i := 0; i <= 10; i++ {
			naturals <- i
		}
		close(naturals)
	}()

	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}
}

func sol_III() {
	random := make(chan int)
	done := make(chan bool)

	go func() {
		defer close(random)

		for i := 0; i < 5; i++ {
			random <- rand.Intn(5)
		}
	}()

	go func() {
		for {
			num, ok := <-random
			if !ok {
				done <- true
			} else {
				fmt.Println(num)
			}
		}
	}()

	<-done
	close(done)
}
