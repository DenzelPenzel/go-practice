/*
Given a string s, find the longest palindromic subsequence's length in s.

A subsequence is a sequence that can be derived from another sequence by deleting some or
no elements without changing the order of the remaining elements.

Example 1:
	Input: s = "bbbab"
	Output: 4
	Explanation: One possible longest palindromic subsequence is "bbbb".

Example 2:
	Input: s = "cbbd"
	Output: 2
	Explanation: One possible longest palindromic subsequence is "bb".

Constraints:
	1 <= s.length <= 1000
	s consists only of lowercase English letters.

*/

package main

func longestPalindromeSubseq(A string) int {
	n := len(A)
	var dp [][]int
	for i := 0; i < n; i++ {
		dp = append(dp, []int{})
		for j := 0; j < n; j++ {
			dp[i] = append(dp[i], 0)
		}
	}
	for i := 0; i < n; i++ {
		dp[i][i] = 1
	}
	for length := 1; length < n; length++ {
		for i := n - length - 1; i >= 0; i-- {
			j := i + length
			if A[i] == A[j] {
				dp[i][j] = dp[i+1][j-1] + 2
			} else {
				dp[i][j] = max(dp[i+1][j], dp[i][j-1])
			}
		}
	}

	return dp[0][n-1]
}

func max(a int, rest ...int) int {
	res := a
	for _, v := range rest {
		if res < v {
			res = v
		}
	}
	return res
}
