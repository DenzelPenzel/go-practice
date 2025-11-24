package main

import "time"

func worker() chan int {
	ch := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()

	return ch
}

func main() {
	timeStart := time.Now()

	// immediately blocks waiting value
	// only after the first receive completes (≈3 s later)
	// it evaluate the second worker() and block on that result, leading to ~6 s total
	//_, _ = <-worker(), <-worker()

	// 1. Start the first worker (instant)
	// This spawns the background goroutine immediately.
	first := worker()

	// 2. Start the second worker (instant)
	// Now both goroutines are sleeping in the background at the same time.
	second := worker()

	// 3. Wait for results
	// Since they started at the same time, they will finish at roughly the same time
	<-first
	<-second

	println(int(time.Since(timeStart).Seconds())) // что выведет - 3 или 6?
}
