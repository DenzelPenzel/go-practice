package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	ch := make(chan int, 5)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch <- rand.Intn(5)
		}
		close(ch)
	}()

	val := <-ch

	fmt.Println(val)

	wg.Wait()
}
