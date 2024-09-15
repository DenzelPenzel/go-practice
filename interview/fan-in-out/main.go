package main

import (
	"fmt"
	"sync"
)

func factorialWorker(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		res := 1
		select {
		case v, ok := <-in:
			if !ok {
				return
			}
			for i := 2; i <= v; i++ {
				res *= i
			}
			out <- res
		}
	}

	// for v := range in {
	// 	result := 1
	// 	for i := 2; i <= v; i++ {
	// 		result *= i
	// 	}
	// 	out <- result
	// }
}

func main() {
	fmt.Println("Hello")

	numWorkers := 5

	numChan := make(chan int, 10)
	resultChan := make(chan int, 10)
	var generationWG sync.WaitGroup

	// Generate numbers in separate goroutines
	go func() {
		for i := 0; i < 10; i++ {
			generationWG.Add(1)
			go func(mm chan<- int, i int, group *sync.WaitGroup) {
				defer group.Done()
				mm <- i + 1
			}(numChan, i, &generationWG)
		}
		generationWG.Wait()
		close(numChan)
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

	// fan-in
	for v := range resultChan {
		fmt.Println("Result:", v)
	}
}
