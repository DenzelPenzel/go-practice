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

func sumSlice() {
	var numbers []*int

	for _, val := range []int{1, 2, 3, 4} {
		x := 5
		numbers = append(numbers, &val, &x)
	}

	numbers = double(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number)
	}

	fmt.Println("=====================")
}

func changePointer() {
	v := 5
	p := &v // store the address of v in p

	fmt.Println(*p) // 5

	changePointerValue(p)

	fmt.Println(*p) // 10
}

// Go passes arguments by value, but in this case the value being copied is the address of v
func changePointerValue(p *int) {
	fmt.Println(p)  // address of v
	fmt.Println(*p) // 5
	*p = 10
}

func main() {
	sumSlice()

	changePointer()
}
