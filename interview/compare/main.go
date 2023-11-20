package main

import "fmt"

func main() {
	// Arrays compare
	var arr1 = [4]int{1, 2, 3, 4}
	var arr2 = [4]int{1, 2, 3, 4}
	fmt.Println(arr1 == arr2) // true

	// Arrays compare
	var str1 = [...]string{"1"}
	var str2 = [...]string{"1"}
	fmt.Println(str1 == str2) // true

	// Slices cannot be compared directly because slice can only be compared to nil
	var str3 = []string{"1"}
	var str4 = []string{"1"}
	fmt.Println(str3 == str4) // error

	var arr3 = []int{1, 2, 3, 4}
	var arr4 = []int{1, 2, 3, 4}
	fmt.Println(arr1 == arr2) // error
}
