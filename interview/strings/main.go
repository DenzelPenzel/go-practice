/*

Develop an algorithm to ascertain whether every character in a string is unique.
In this task, the use of additional storage structures is not permitted.
Given a string, your task is to return a boolean value.
True indicates that all characters are distinct, while false implies the presence of duplicate characters.

It's important to ensure that the characters in the string are limited to ASCII characters.
The length of the string will not exceed 3000 characters
*/

package main

import (
	"fmt"
	"strings"
)

func isUniqueString(s string) bool {
	for i, v := range s {
		if strings.Index(s, string(v)) != i {
			return false
		}
	}
	return true
}

func isUniqString_II(s string) bool {
	if len(s) == 0 || len(s) > 3000 {
		return false
	}
	var mark1, mark2, mark3, mark4 uint64
	var mark *uint64
	for _, r := range s {
		x := uint64(r)
		if x < 64 {
			mark = &mark1
		} else if x < 128 {
			mark = &mark2
			x -= 64
		} else if x < 192 {
			mark = &mark3
			x -= 128
		} else {
			mark = &mark4
			x -= 192
		}
		if (*mark)&(1<<x) != 0 {
			return false
		}
		*mark = (*mark) | uint64(1<<x)
	}
	return true
}

func main() {
	res := isUniqueString("abca")
	res2 := isUniqString_II("abca")

	fmt.Println("pass 'abca' want return false, got", res)
	fmt.Println("pass 'abca' want return false, got", res2)
}
