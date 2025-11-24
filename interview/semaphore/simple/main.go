package main

import "fmt"

type sema chan struct{}

func New(size int) sema {
	return make(sema, size)
}

func (s sema) Acquire(k int) {
	for i := 0; i < k; i++ {
		s <- struct{}{}
	}
}

func (s sema) Release(k int) {
	for i := 0; i < k; i++ {
		<-s
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	sem := New(len(nums))

	for _, num := range nums {
		go func(num int) {
			fmt.Println(num)
			sem.Acquire(1)
		}(num)
	}

	sem.Release(len(nums))
}
