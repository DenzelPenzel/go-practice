package main

import "fmt"

type Student struct {
	Name string
}

// Return print results
func main() {
	// Pointer types - compare pointer addresses
	fmt.Println(&Student{Name: "vasya"} == &Student{Name: "vasya"}) // false

	// Non-pointer types - compare the value of each attribute
	fmt.Println(Student{Name: "vasya"} == Student{Name: "vasya"}) // true
}
