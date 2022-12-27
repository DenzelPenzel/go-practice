package main

/*
You are given a string s that consists of lower case English letters and brackets.

Reverse the strings in each pair of matching parentheses, starting from the innermost one.

Your result should not contain any brackets.

Example 1:
	Input: s = "(abcd)"
	Output: "dcba"

Example 2:
	Input: s = "(u(love)i)"
	Output: "iloveu"
	Explanation: The substring "love" is reversed first, then the whole string is reversed.

Example 3:
	Input: s = "(ed(et(oc))el)"
	Output: "leetcode"
	Explanation: First, we reverse the substring "oc", then "etco", and finally, the whole string.

Constraints:
	1 <= s.length <= 2000
	s only contains lower case English characters and parentheses.
	It is guaranteed that all parentheses are balanced.
*/

func reverseParentheses(s string) string {
	st := make([]string, 0)
	current := ""

	for _, x := range s {
		if string(x) == ")" {
			for st[len(st)-1] != "(" {
				ch := st[len(st)-1]
				st = st[:len(st)-1]
				current = ch + current
			}
			st = st[:len(st)-1]
			st = append(st, reverse(current))
			current = ""
		} else if string(x) == "(" {
			if len(current) > 0 {
				st = append(st, current)
			}
			st = append(st, "(")
			current = ""
		} else {
			current = current + string(x)
		}
	}

	if len(current) > 0 {
		st = append(st, current)
	}

	res := ""

	for len(st) > 0 {
		ch := st[len(st)-1]
		st = st[:len(st)-1]
		res = ch + res
	}

	return res
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
