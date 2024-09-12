package main

import "fmt"

func double(numbers []*int) []*int {
	res := 0
	for i := 0; i < len(numbers); i++ {
		*numbers[i] = *numbers[i] * 2
		res += *numbers[i]
	}
	numbers = append(numbers, &res)
	return numbers
}

func main() {
	var numbers []*int

	for _, val := range []int{1, 2, 3, 4} {
		x := 5
		numbers = append(numbers, &val, &x)
	}

	numbers = double(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number)
	}
}
