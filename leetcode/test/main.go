/*
You are given a string s consisting of the characters 'a', 'b', and 'c' and a non-negative integer k.
Each minute, you may take either the leftmost character of s, or the rightmost character of s.

Return the minimum number of minutes needed for you to take at least k of each character,
or return -1 if it is not possible to take k of each character.

Example 1:
Input: s = "aabaaaacaabc", k = 2
Output: 8
Explanation:
Take three characters from the left of s. You now have two 'a' characters, and one 'b' character.
Take five characters from the right of s. You now have four 'a' characters, two 'b' characters, and two 'c' characters.
A total of 3 + 5 = 8 minutes is needed.
It can be proven that 8 is the minimum number of minutes needed.

Example 2:
Input: s = "a", k = 1
Output: -1
Explanation: It is not possible to take one 'b' or 'c' so return -1.

Constraints:
1 <= s.length <= 105
s consists of only the letters 'a', 'b', and 'c'.
0 <= k <= s.length
*/

package main

import "math"

func takeCharacters(A string, k int) int {
	n := len(A)
	A = A + A
	start, end := 0, 0
	freq := make([]int, 3)
	minWindow := math.MaxInt32

	for end < len(A) {
		freq[A[end]-'a']++
		end++
		for end-start <= n && min(freq) >= k && (start == 0 || start <= len(A)/2 && end >= len(A)/2) {
			if end-start < minWindow {
				minWindow = end - start
			}
			freq[A[start]-'a']--
			start++
		}
	}
	if minWindow <= n {
		return minWindow
	}
	return -1
}

func min(array []int) int {
	minVal := array[0]
	for _, val := range array {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}
