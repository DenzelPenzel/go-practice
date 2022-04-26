/*
Given an array nums of integers,
we need to find the maximum possible sum of elements of the array such
that it is divisible by three.

Example 1:
	Input: nums = [3,6,5,1,8]
	Output: 18
	Explanation: Pick numbers 3, 6, 1 and 8 their sum is 18 (maximum sum divisible by 3).

Example 2:
	Input: nums = [4]
	Output: 0
	Explanation: Since 4 is not divisible by 3, do not pick any number.

Example 3:
	Input: nums = [1,2,3,4,4]
	Output: 12
	Explanation: Pick numbers 1, 3, 4 and 4 their sum is 12 (maximum sum divisible by 3).

Constraints:
	1 <= nums.length <= 4 * 10^4
	1 <= nums[i] <= 10^4

*/

package leetcode

func maxSumDivThree(nums []int) int {
	n := len(nums)
	dp := [][]int{}

	for i := 0; i < n; i++ {
		dp = append(dp, []int{0, 0, 0})
	}

	dp[0][nums[0]%3] = nums[0]

	for i := 1; i < n; i++ {
		for k := 0; k < 3; k++ {
			mod := (dp[i-1][k] + nums[i]) % 3
			dp[i][mod] = Max(dp[i][mod], dp[i-1][k]+nums[i])
			dp[i][k] = Max(dp[i][k], dp[i-1][k])
		}
	}

	res := 0

	for i := 0; i < n; i++ {
		res = Max(res, dp[i][0])
	}

	return res
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
