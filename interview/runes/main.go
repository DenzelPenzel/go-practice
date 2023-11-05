package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Boat 🚢."
	fmt.Println("length of string?", s, len(s))

	length := utf8.RuneCountInString(s)

	fmt.Println("length of string?", length)

	for i, c := range s {
		// rune
		fmt.Printf("Position %d of '%s'\n", i, string(c))
	}
}
