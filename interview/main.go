package main

import (
	"fmt"
	"sync"
)

func factorialWorker(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		result := 1
		for i := 2; i <= v; i++ {
			result *= i
		}
		out <- result
	}
}

func main() {
	numWorkers := 2
	numChan := make(chan int, 10)
	resultChan := make(chan int, 10)
	var generationWG sync.WaitGroup

	go func() {
		defer close(numChan)

		for i := 0; i < 10; i++ {
			generationWG.Add(1)

			go func(i int, group *sync.WaitGroup) {
				defer group.Done()
				numChan <- i + 1
			}(i, &generationWG)
		}

		generationWG.Wait()
	}()

	var wg sync.WaitGroup

	// Start the factorial workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go factorialWorker(numChan, resultChan, &wg)
	}

	// Wait for all workers to finish and then close resultChan
	go func() {
		wg.Wait()
		close(resultChan) // <- if comment this line will be deadlock!
	}()

	// keep reading from resultChan until it is closed
	for v := range resultChan {
		fmt.Println("Result:", v)
	}
}
