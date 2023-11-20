package main

import "fmt"

type A struct {
	s string
}

func foo(s string) *A {
	a := new(A) // escaped, heap allocation
	a.s = s
	return a
}

func main() {
	a := foo("hello")
	b := a.s + " world" // stack
	c := b + "!!!!"     // heap, escape will occur through fmt.Println(a ...interface{})
	fmt.Println(c)
}
