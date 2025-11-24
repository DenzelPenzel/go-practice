package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

func sol1() {
	const n = 1000
	var counter int32

	incr := make(chan int, n)

	for i := 0; i < 1000; i++ {
		go func() {
			incr <- 1
			if atomic.AddInt32(&counter, 1) == int32(n) {
				close(incr)
			}
		}()
	}

	var res int

	for num := range incr {
		res += num
	}

	fmt.Println(res)
}

func sol2() {
	runtime.GOMAXPROCS(1)

	const n = 1000
	var counter int64

	for i := 0; i < n; i++ {
		go func() {
			atomic.AddInt64(&counter, 1)
		}()
	}

	for atomic.LoadInt64(&counter) != int64(n) {
		runtime.Gosched()
	}

	fmt.Println(counter)
}

func sol3() {
	runtime.GOMAXPROCS(1)

	const n = 1000
	incr := make(chan int, n)

	for i := 0; i < n; i++ {
		go func() {
			incr <- 1
		}()
	}

	res := 0
	for i := 0; i < n; i++ {
		res += <-incr
	}
	fmt.Println(res)
}

func main() {
	sol1()
	sol2()
	sol3()
}
