package main

import "fmt"

// Return print results
func main() {
	test1()
	test2()
}

// ==========================================================================================================
// copy slice
func test1() {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:] // copy operation

	str2[1] = "new"
	fmt.Println(str1) // [a b new]

	str2 = append(str2, "z", "x", "y") // this operation involve to expand and generate a new array, not effect str1
	fmt.Println(str1)                  // [a b new]
	fmt.Println(str2)                  // [b new z x y]
}

// ==========================================================================================================
// extend slice
func test2() {
	fmt.Println("\n=============")

	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	fmt.Println(s) // 0, 0, 0, 0, 0, 1, 2, 3
}
