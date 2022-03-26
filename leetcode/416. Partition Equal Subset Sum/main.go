package leetcode

func canPartition(nums []int) bool {
	total := 0
	for _, v := range nums {
		total += v
	}

	if total%2 != 0 {
		return false
	}

	target := total >> 1

	dp := make([]bool, target+1)

	dp[0] = true

	for _, coin := range nums {
		for sum := target; sum >= coin; sum-- {
			dp[sum] = dp[sum] || dp[sum-coin]
		}
	}

	return dp[target]
}
