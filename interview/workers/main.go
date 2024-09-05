package main

import "fmt"

const numJobs = 5

func worker(id int, f func(int) int, jobs <-chan int, res chan<- int) {
	for j := range jobs {
		res <- f(j)
	}
}

/*
Task to split the processes into several G without creating a new goroutine each time:

 1. To do this, we will create a channel with jobs and a resulting channel

 2. For each worker, we will create a goroutine that will wait for a new job,
    apply the specified function to it and publish the response in the resulting channel
*/
func sol() {
	jobs := make(chan int, numJobs)
	res := make(chan int, numJobs)

	multiplier := func(x int) int {
		return x * 10
	}

	for w := 1; w <= 3; w++ {
		go worker(w, multiplier, jobs, res)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}

	close(jobs)

	for i := 1; i <= numJobs; i++ {
		fmt.Println(<-res)
	}
}
