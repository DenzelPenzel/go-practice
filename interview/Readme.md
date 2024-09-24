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


### Sync G work
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

### Fan-Out, Fan-In

```
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
	numWorkers := 5

	numChan := make(chan int, 10)
	resultChan := make(chan int, 10)
	var generationWG sync.WaitGroup

	// Generate numbers in separate goroutines
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

	// fan-in
	for v := range resultChan {
		fmt.Println("Result:", v)
	}
}
```


### Pipeline
```
func gen(words ...string) <-chan string {
	out := make(chan string)

	go func() {
		for _, v := range words {
			out <- v
		}
		close(out)
	}()
	return out
}

func toUpper(in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		for s := range in {
			out <- strings.ToUpper(s)
		}
		close(out)
	}()
	return out
}

func exclaim(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for s := range in {
			out <- s + "!"
		}
		close(out)
	}()
	return out
}

func main() {
	pipeline := exclaim(toUpper(gen("hello", "world", "golang", "pipeline")))

	for v := range pipeline {
		fmt.Println(v)
	}
}
```


### Rate limit
```
func main() {
    // Create a ticker that ticks every 500ms
    limiter := time.Tick(500 * time.Millisecond)

    jobs := []int{1, 2, 3, 4, 5}

    for _, job := range jobs {
        <-limiter // Wait for the next tick
        go func(j int) {
            fmt.Printf("Processing job %d at %v\n", j, time.Now())
        }(job)
    }

    // Wait to allow goroutines to finish
    time.Sleep(3 * time.Second)
}
```


### Publish/Subscribe (Pub/Sub)
```
// Broker manages subscriptions and publishing
type Broker struct {
    subscribers map[chan string]struct{}
    lock        sync.RWMutex
}

// NewBroker creates a new Broker
func NewBroker() *Broker {
    return &Broker{
        subscribers: make(map[chan string]struct{}),
    }
}

// Subscribe returns a channel to receive messages
func (b *Broker) Subscribe() <-chan string {
    ch := make(chan string)
    b.lock.Lock()
    b.subscribers[ch] = struct{}{}
    b.lock.Unlock()
    return ch
}

// Unsubscribe removes a subscriber
func (b *Broker) Unsubscribe(ch <-chan string) {
    b.lock.Lock()
    delete(b.subscribers, ch.(chan string))
    close(ch.(chan string))
    b.lock.Unlock()
}

// Publish sends a message to all subscribers
func (b *Broker) Publish(msg string) {
    b.lock.RLock()
    defer b.lock.RUnlock()
    for ch := range b.subscribers {
        // Non-blocking send
        select {
        case ch <- msg:
        default:
            // If subscriber is not ready, skip
        }
    }
}

func main() {
    broker := NewBroker()

    // Subscriber 1
    sub1 := broker.Subscribe()
    go func() {
        for msg := range sub1 {
            fmt.Println("Subscriber 1 received:", msg)
        }
    }()

    // Subscriber 2
    sub2 := broker.Subscribe()
    go func() {
        for msg := range sub2 {
            fmt.Println("Subscriber 2 received:", msg)
        }
    }()

    // Publish messages
    messages := []string{"Hello", "World", "Golang", "Concurrency"}

    for _, msg := range messages {
        broker.Publish(msg)
    }

    // Unsubscribe
    broker.Unsubscribe(sub1)
    broker.Unsubscribe(sub2)

    // Wait to ensure all messages are processed
    // In real applications, use synchronization
    // Here, sleep for simplicity
    time.Sleep(time.Second)
}
```
