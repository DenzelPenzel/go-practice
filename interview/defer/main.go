package main

import "fmt"

func main() {
	fmt.Println(test1()) // 11

	test3()
	// test4()
	test5()
}

// ==========================================================================================================

func test1() (result int) {
	defer func() {
		result++
	}()

	return 10
}

// ==========================================================================================================

func test3() {
	fmt.Println("\n=============")

	var i1 int = 10
	var k = 20
	var i2 *int = &k

	defer printInt("i1", i1)                   // 10
	defer printInt("i2 as value", *i2)         // 20 bcs we copy value
	defer printIntPointer("i2 as pointer", i2) // 2020 bcs we copy reference

	i1 = 1010
	*i2 = 2020

}

// ==========================================================================================================

func test4() {
	fmt.Println("\n=============")

	// when called runtime.deferreturn, it will be popped from the stack and executed in the reverse order
	defer func() {
		fmt.Println("Hello 1\n")
	}()
	defer func() {
		fmt.Printf("Hello 2\n")
	}()
	defer func() {
		fmt.Println("Hello 3\n")
	}()

	// panic does not terminate defer execution
	panic("Error")
}

// ==========================================================================================================

func test5() {
	fmt.Println("\n=============")

	var calc func(idx string, a, b int) int
	calc = func(idx string, a, b int) int {
		res := a + b
		fmt.Println(idx, a, b, res)
		return res
	}

	a := 1
	b := 2

	// parameters of the calling function, will be calculated when defining
	// calc("10", a, b) and calc("20", a, b) - call immediately
	// calc("2", ...)   and calc("1"...) - in reverse order
	defer calc("1", a, calc("10", a, b))
	a = 0

	defer calc("2", a, calc("20", a, b))
	b = 1
}

func changeI(i *int) {
	*i = 10
}

func printInt(v string, i int) {
	fmt.Printf("%s=%d\n", v, i)
}

func printIntPointer(v string, i *int) {
	fmt.Printf("%s=%d\n", v, *i)
}
