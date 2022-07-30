package leetcode

import (
	"math"
)

func palindromePartition(s string, k int) int {
	dp := [][]int{}
	memo := [][]int{}
	n := len(s)

	for i := 0; i < n; i++ {
		dp = append(dp, []int{})
		memo = append(memo, []int{})
		for j := 0; j < n+1; j++ {
			dp[i] = append(dp[i], 0)
			memo[i] = append(memo[i], -1)
		}
	}

	for length := 1; length < n; length++ {
		for i := 0; i < n-length; i++ {
			j := i + length
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1]
			} else {
				dp[i][j] = dp[i+1][j-1] + 1
			}
		}
	}

	return DFS(s, dp, &memo, 0, k)
}

func DFS(s string, dp [][]int, memo *[][]int, i, k int) int {
	if i >= len(s) {
		if k == 0 {
			return 0
		}
		return int(math.MaxInt32)
	}

	if k <= 0 {
		return int(math.MaxInt32)
	}

	if (*memo)[i][k] != -1 {
		return (*memo)[i][k]
	}

	res := int(math.MaxInt32)

	for j := i; j < len(s); j++ {
		res = min(res, dp[i][j]+DFS(s, dp, memo, j+1, k-1))
	}

	(*memo)[i][k] = res

	return res
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
