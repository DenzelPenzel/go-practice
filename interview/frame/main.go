package main

func main() {
	count := 10

	println("count: \tValue of [", count, "]\tAddr Of [", &count, "]")

	// Pass the "value of" the count
	increamentByValue(count)

	// Share the "address of" the count but it's still a "pass by value"
	// the only difference is the value you are passing is an address instead of an integer
	increamentByRef(&count)

	println("count:\tValue Of[", count, "]\tAddr Of[", &count, "]")
}

// Goroutine needs to frame a new section of memory on the stack
func increamentByValue(inc int) {
	// value will be copied, passed and placed into the new frame for the increment func
	inc++
	println("inc by val: \tValue of [", inc, "]\tAddr of [", &inc, "]")
}

func increamentByRef(inc *int) {
	*inc++
	// inc != &inc
	println("inc ref:\tValue Of[", inc, "]\tAddr Of reference[", &inc, "]\tValue Points To[", *inc, "]")
}
