package main

import "fmt"

// There are two channels
// Numbers are written into the first
// The numbers need to be read from the first channel and squared
// Result needs to be written into the second channel
func main() {
	nums := make(chan int)
	squares := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			nums <- i
		}
		close(nums)
	}()

	go func() {
		for num := range nums {
			squares <- num * num
		}
		close(squares)
	}()

	for num := range squares {
		fmt.Println(num)
	}
}
