/*
A string is called a happy prefix if is a non-empty prefix which is also a suffix (excluding itself).

Given a string s, return the longest happy prefix of s.

Return an empty string "" if no such prefix exists.

Example 1:
    Input: s = "level"
    Output: "l"
    Explanation:
        s contains 4 prefix excluding itself ("l", "le", "lev", "leve"), and suffix ("l", "el", "vel", "evel").
        The largest prefix which is also suffix is given by "l".

Example 2:
    Input: s = "ababab"
    Output: "abab"
    Explanation:
        "abab" is the largest prefix which is also suffix. They can overlap in the original string.


Constraints:
    1 <= s.length <= 10^5
    s contains only lowercase English letters.
*/

package main

// Use Longest Prefix Suffix (KMP-table) or String Hashing
func longestPrefix(s string) string {
	lps := make([]int, len(s))
	i, j := 0, 1

	for j < len(s) {
		if s[i] == s[j] {
			lps[j] = i + 1
			i++
			j++
		} else {
			if i != 0 {
				i = lps[i-1]
			} else {
				lps[j] = 0
				j++
			}
		}
	}

	return s[0:i]
}
