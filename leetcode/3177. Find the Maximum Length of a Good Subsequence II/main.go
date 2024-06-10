/*
You are given an integer array nums and a non-negative integer k.

A sequence of integers seq is called good if there are
at most k indices i in the range [0, seq.length - 2] such that seq[i] != seq[i + 1].

Return the maximum possible length of a good subsequence of nums.

Example 1:
	Input: nums = [1,2,1,1,3], k = 2
	Output: 4
	Explanation:
		The maximum length subsequence is [1,2,1,1,3].

Example 2:
	Input: nums = [1,2,3,4,5,1], k = 0
	Output: 2
	Explanation:
		The maximum length subsequence is [1,2,3,4,5,1].


Constraints:
	1 <= nums.length <= 5 * 10^3
	1 <= nums[i] <= 10^9
	0 <= k <= min(50, nums.length)
*/

package main

func maximumLength(A []int, k int) int {
	n := len(A)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = 1
		}
	}
	res := 1
	dp2 := make([]map[int]int, k+1)
	dp_k := make([]int, k+1)
	for i := range dp2 {
		dp2[i] = make(map[int]int)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < k+1; j++ {
			// check all possible previous transition counts
			for jj := 0; jj < j; jj++ {
				// update dp by considering previous transition counts
				dp[i][j] = max(dp[i][j], dp_k[jj]+1)
			}
			// check if the current number has been seen before at this transition count
			if idx, ok := dp2[j][A[i]]; ok {
				dp[i][j] = max(dp[i][j], idx+1)
			}
			res = max(res, dp[i][j])
		}
		for j := 0; j < k+1; j++ {
			dp2[j][A[i]] = max(dp2[j][A[i]], dp[i][j])
			dp_k[j] = max(dp_k[j], dp[i][j])
		}
	}
	return res
}
