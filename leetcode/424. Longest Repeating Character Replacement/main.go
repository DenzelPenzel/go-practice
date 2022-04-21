/*
You are given a string s and an integer k.

You can choose any character of the string and
change it to any other uppercase English character.

You can perform this operation at most k times.

Return the length of the longest substring containing the same letter
you can get after performing the above operations.


Example 1:
	Input: s = "ABAB", k = 2
	Output: 4
	Explanation: Replace the two 'A's with two 'B's or vice versa.

Example 2:
	Input: s = "AABABBA", k = 1
	Output: 4
	Explanation: Replace the one 'A' in the middle with 'B' and form "AABBBBA".
	The substring "BBBB" has the longest repeating letters, which is 4.

Constraints:
	1 <= s.length <= 105
	s consists of only uppercase English letters.
	0 <= k <= s.length
*/

package leetcode

func characterReplacement(s string, k int) int {
	res := 0
	start, end := 0, 0
	freq := make([]int, 26)

	for end < len(s) {
		freq[s[end]-'A'] += 1
		end += 1
		max_freq := getMaxFreq(freq)
		for end-start-max_freq > k {
			freq[s[start]-'A'] -= 1
			start += 1
			max_freq = getMaxFreq(freq)
		}
		res = max(res, end-start)
	}
	return res
}

func getMaxFreq(freq []int) int {
	maxFreq := 0
	for i := 0; i < 26; i++ {
		maxFreq = max(maxFreq, freq[i])
	}
	return maxFreq
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
