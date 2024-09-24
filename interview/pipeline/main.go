package main

import (
	"fmt"
	"strings"
)

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
