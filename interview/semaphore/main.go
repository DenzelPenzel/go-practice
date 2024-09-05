package main

import "fmt"

type sem chan struct{}

func New(n int) sem {
	return make(sem, n)
}

func (s sem) Inc(k int) {
	for i := 0; i < k; i++ {
		s <- struct{}{}
	}
}

func (s sem) Dec(k int) {
	for i := 0; i < k; i++ {
		<-s
	}
}

// Make a custom waitGroup on a semaphore

/*
We want to make a semaphore that will wait for five goroutines to complete!

1. Create buffered channel, inside each goroutine we put a value in it
2. At the end we will expect that everything is ok - we will subtract all the values ​​from the channel
*/
func main() {
	nums := []int{1, 2, 3, 4, 5}
	n := len(nums)

	s := New(n)

	for _, x := range nums {
		go func(x int) {
			fmt.Println(x)
			s.Inc(1)
		}(x)
	}

	s.Dec(n)
}
