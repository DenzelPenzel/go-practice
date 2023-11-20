package main

import "fmt"

// Find the problem with the code below

type query func(string) string

// The above code has a serious memory leak problem
// go fn(i)4 coroutines will actually be started after the code is executed
func exec(name string, rest ...query) string {
	// because ch is non-buffered, only one coroutine may be written successfully
	// the other three coroutines will always wait for writing in the background
	ch := make(chan string)
	fn := func(i int) {
		ch <- rest[i](name)
	}
	for i, _ := range rest {
		go fn(i)
	}
	return <-ch
}

func main() {
	res := exec("run run run", func(s string) string {
		return s + " func1"
	}, func(s string) string {
		return s + " func2"
	}, func(s string) string {
		return s + " func3"
	}, func(s string) string {
		return s + " func4"
	})

	fmt.Println(res)
}
