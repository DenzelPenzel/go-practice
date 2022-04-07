package leetcode

func numTrees(n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1
	for i := 2; i <= n; i++ {
		for k := 1; k <= i; k++ {
			dp[i] += dp[k-1] * dp[i-k]
		}
	}
	return dp[n]
}
