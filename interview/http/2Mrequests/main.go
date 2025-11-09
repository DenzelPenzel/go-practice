package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
)

const (
	TotalRequests  = 2000000
	RequestURL     = "https://example.com/api"
	TimeoutSeconds = 10
)

func worker(id int, jobs <-chan int, results chan<- error, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		results <- sendRequest(client, j)
	}
}

func main() {
	workers := runtime.NumCPU()
	jobs := make(chan int, workers)
	results := make(chan error, TotalRequests)

	client := &http.Client{
		Timeout: TimeoutSeconds * 1e9,
	}

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, client, &wg)
	}

	go func() {
		for i := 0; i < TotalRequests; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()
	close(results)

	var success, failure int
	for err := range results {
		if err != nil {
			failure++
			log.Println("Request failed:", err)
		} else {
			success++
		}
	}

	fmt.Printf("Total Requests: %d\n", TotalRequests)
	fmt.Printf("Successful Requests: %d\n", success)
	fmt.Printf("Failed Requests: %d\n", failure)
}

func sendRequest(client *http.Client, jobID int) error {
	req, err := http.NewRequest("GET", RequestURL, nil)
	if err != nil {
		return fmt.Errorf("job %d: creating request failed: %w", jobID, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("job %d: request failed: %w", jobID, err)
	}

	defer resp.Body.Close()

	_, readErr := io.Copy(io.Discard, resp.Body)
	if readErr != nil {
		return fmt.Errorf("job %d: reading response failed: %w", jobID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("job %d: unexpected status code: %d", jobID, resp.StatusCode)
	}

	return nil
}
