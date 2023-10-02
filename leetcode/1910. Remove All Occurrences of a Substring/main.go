/*
Given two strings s and part, perform the following operation on s until all occurrences of the substring part are removed:

Find the leftmost occurrence of the substring part and remove it from s.
Return s after removing all occurrences of part.

A substring is a contiguous sequence of characters in a string.

Example 1:

	Input: s = "daabcbaabcbc", part = "abc"
	Output: "dab"
	Explanation:
	    The following operations are done:
	    - s = "daabcbaabcbc", remove "abc" starting at index 2, so s = "dabaabcbc".
	    - s = "dabaabcbc", remove "abc" starting at index 4, so s = "dababc".
	    - s = "dababc", remove "abc" starting at index 3, so s = "dab".
	    Now s has no occurrences of "abc".

Example 2:

	Input: s = "axxxxyyyyb", part = "xy"
	Output: "ab"
	Explanation:
	    The following operations are done:
	    - s = "axxxxyyyyb", remove "xy" starting at index 4 so s = "axxxyyyb".
	    - s = "axxxyyyb", remove "xy" starting at index 3 so s = "axxyyb".
	    - s = "axxyyb", remove "xy" starting at index 2 so s = "axyb".
	    - s = "axyb", remove "xy" starting at index 1 so s = "ab".
	    Now s has no occurrences of "xy".

Constraints:

	1 <= s.length <= 1000
	1 <= part.length <= 1000
	s​​​​​​ and part consists of lowercase English letters.
*/

package p1910removealloccurrencesofsubstring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type Item struct {
	char         rune
	patternIndex int
}

// Naive solution using the stack
func removeOccurrences_II(s string, pattern string) string {
	stack := make([]rune, 0)
	for _, ch := range s {
		stack = append(stack, ch)
		if len(stack) >= len(pattern) {
			if string(stack[(len(stack)-len(pattern)):]) == pattern {
				stack = stack[:len(stack)-len(pattern)]
			}
		}
	}
	return string(stack)
}

// Optimal solution
func removeOccurrences(s string, pattern string) string {
	// build longest prefix suffix table
	lps := buildLPS(pattern)

	stack := make([]Item, 0)
	stack = append(stack, Item{char: ' ', patternIndex: 0})

	for _, ch := range s {
		p := stack[len(stack)-1]
		k := p.patternIndex

		for k > 0 && ch != rune(pattern[k]) {
			k = lps[k-1]
		}

		if ch == rune(pattern[k]) {
			k++
		}

		stack = append(stack, Item{char: ch, patternIndex: k})

		if stack[len(stack)-1].patternIndex == len(pattern) {
			stack = stack[:len(stack)-len(pattern)]
		}
	}

	res := ""

	for len(stack) > 1 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = string(p.char) + res
	}

	return res
}

func buildLPS(s string) []int {
	var lps []int
	lps = append(lps, 0)
	k := 0
	for i := 1; i < len(s); i++ {
		for k > 0 && s[i] != s[k] {
			k = lps[k-1]
		}
		if s[i] == s[k] {
			k++
		}
		lps = append(lps, k)
	}
	return lps
}

func Test_removeOccurrences(t *testing.T) {
	for _, tc := range []struct {
		s       string
		pattern string
		want    string
	}{
		{"daabcbaabcbc", "abc", "dab"},
		{"axxxxyyyyb", "xy", "ab"},
	} {
		t.Run(
			fmt.Sprintf("%+v", tc.s), func(t *testing.T) {
				require.Equal(t, tc.want, removeOccurrences(tc.s, tc.pattern))
			})
	}
}
