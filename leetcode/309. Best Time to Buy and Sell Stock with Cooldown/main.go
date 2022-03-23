package leetcode

import "math"

func maxProfit(A []int) int {
	n := len(A)
	dp := make([][2]int, n)

	dp[0][1] = -A[0] // 0 - sell, 1 - buy
	dp[0][0] = -math.MaxInt32

	for i := 1; i < n; i++ {
		dp[i][0] = Max(dp[i-1][0], dp[i-1][1]+A[i]) // sell

		if i-2 >= 0 {
			dp[i][1] = Max(dp[i-1][1], Max(dp[i-2][0]-A[i], -A[i])) // buy
		} else {
			dp[i][1] = Max(dp[i-1][1], -A[i])
		}
	}

	return Max(dp[n-1][0], 0)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
