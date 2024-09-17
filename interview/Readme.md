### Worker pool 

- Task to split the processes into several G without creating a new goroutine each time
  
```
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

```


### Close channel 
```
func main() {
	var wc sync.WaitGroup
	m := make(chan string, 3)
	fff := sync.Mutex{}

	go func() {
		defer close(m)
		for i := 0; i < 5; i++ {
			wc.Add(1)
			go func(mm chan<- string, i int, group *sync.WaitGroup) {
				defer wc.Done()
				fff.Lock() // do you really need it here?
				mm <- fmt.Sprintf("Goroutine %s", strconv.Itoa(i))
				fff.Unlock()
			}(m, i, &wc)
		}
		wc.Wait()
	}()

	for q := range m {
		fmt.Println(q)
	}

	// for {
	// 	select {
	// 	case q, ok := <-m:
	// 		if !ok {
	// 			return
	// 		} else {
	// 			fmt.Println(q)
	// 		}
	// 	}
	// }
}
```


### Use channels to coordinate access to shared state
```
func main() {
	counterChan := make(chan int)
	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	// Counter goroutine that listens for increments
	go func() {
		counter := 0
		for inc := range counterChan {
			counter += inc
		}
		fmt.Println("Final Counter:", counter) // This will print after the channel is closed
		doneChan <- struct{}{}
	}()

	// Launch 1000 goroutines to send increment requests
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counterChan <- 1 // Send increment
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the counter channel to signal completion
	close(counterChan)

	<-doneChan
}
```