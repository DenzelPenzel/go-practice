package main

import (
	"fmt"
	"strings"
	"sync"
)

/*
Task to split the processes into several G without creating a new goroutine each time:

 1. To do this, we will create a channel with jobs and a resulting channel

 2. For each worker, we will create a goroutine that will wait for a new job,
    apply the specified function to it and publish the response in the resulting channel
*/

const workerCount = 5

// Worker Pool Pattern
// This pattern utilizes a fixed number of worker goroutines that process tasks from a shared job queue
func main() {
	var wg sync.WaitGroup
	jobQueue := make(chan string)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(&wg, jobQueue)
	}

	data := getUrls()
	for _, v := range data {
		jobQueue <- v
	}

	close(jobQueue) // Close the job queue once all jobs are submitted

	wg.Wait()
}

// Each worker processes jobs concurrently
func worker(wg *sync.WaitGroup, jobQueue <-chan string) {
	defer wg.Done()
	for url := range jobQueue {
		if checkDomain(url) {
			fmt.Printf("%s found \n", url)
		}
	}
}

func getUrls() []string {
	return []string{"1.test", "2.test", "abc.com", "1000.test"}
}

func checkDomain(host string) bool {
	return strings.Contains(host, ".test")
}
