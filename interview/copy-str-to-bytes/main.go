package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ==========================================================================================================
/*
Convert the string into a byte array

Converting a string into a slice will produce a copy.
Strictly speaking, a memory copy will occur whenever a type force transfer occurs. Then here comes the problem.

Frequent memory copy operations sound unfriendly to performance. Is there any way to convert a string into a slice without copying it?

Solution:
	StringHeader It is the underlying structure of string in go.

	type StringHeader struct {
		Data uintptr
		Len  int
	}

	SliceHeader it is the underlying structure of slices in go.

	type SliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}

	Then if you want to convert the two at the bottom level, you only need to force the address of StringHeader into SliceHeader.
	Then go has a very powerful package called unsafe.

	- unsafe.Pointer(&s) - method can get the address of variable s

	- (*reflect.StringHeader)(unsafe.Pointer(&a)) - method can convert string s into the form of the underlying structure

	- (*[]byte)(unsafe.Pointer(&ss)) - method can convert the ss underlying structure into a pointer to a byte slice

	- then by *converting it to the actual content pointed to by the pointer
*/

func main() {
	s := "abcd"

	ss := *(*reflect.StringHeader)(unsafe.Pointer(&s))

	fmt.Printf("%v\n", ss)

	b := *(*[]byte)(unsafe.Pointer(&ss))
	fmt.Printf("%v", b)

}
