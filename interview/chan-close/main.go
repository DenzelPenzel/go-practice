package main

import (
	"fmt"
	"strconv"
	"sync"
)

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
