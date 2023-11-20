package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	runtime.GOMAXPROCS(1)

	test1()
	test2()
}

func test1() {
	wg := sync.WaitGroup{}
	wg.Add(20)

	// output determines which Goroutine is prioritized by the scheduler.
	// therefore first output is the last created goroutine j = 9
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("i: ", i) // return 10 each time
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("j: ", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func test2() {
	fmt.Println("\n=============")

	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)

	int_chan <- 1
	string_chan <- "hello world"

	// the result is the random execution due the fact that Go fairly select one for execution when multiple reads are available
	select {
	case val := <-int_chan:
		fmt.Println(val)
	case val := <-string_chan:
		fmt.Println(val)
	}
}
