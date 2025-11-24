package main

import (
	"fmt"
	"sync"
)

func worker(id int, fn func(int) int, jobs <-chan int, res chan<- int) {
	fmt.Printf("Worker %d started\n", id)
	for j := range jobs {
		res <- fn(j)
	}
}

func main() {
	const (
		workerCount = 3
		numJobs     = 10
	)
	jobs := make(chan int, numJobs)
	res := make(chan int, numJobs)

	multi := func(x int) int {
		return x * x
	}

	var wg sync.WaitGroup

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			worker(id, multi, jobs, res)
		}(w, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(res)
	}()

	for num := range res {
		fmt.Println(num)
	}
}
