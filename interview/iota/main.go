package main

import "fmt"

const (
	a = iota
	b = iota
)

const (
	name = "hello"
	c    = iota
	d    = iota
)

func main() {
	fmt.Println(a) // 0
	fmt.Println(b) // 1
	fmt.Println(c) // 1
	fmt.Println(d) // 2
}
