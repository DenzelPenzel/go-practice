package leetcode

func maxPoints(A [][]int) int64 {
	n := len(A[0])
	dp := make([]int, n)

	for _, row := range A {
		for i := 1; i < n; i++ {
			dp[i] = Max(dp[i], dp[i-1]-1)
		}

		for i := n - 2; i >= 0; i-- {
			dp[i] = Max(dp[i], dp[i+1]-1)
		}

		for i, v := range row {
			dp[i] += v
		}

	}

	res := 0

	for i := 0; i < n; i++ {
		res = Max(res, dp[i])
	}

	return int64(res)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
