/*
Given a non-empty array nums containing only positive integers,
find if the array can be partitioned into two subsets such that the sum of elements in both subsets is equal.

Example 1:
    Input: nums = [1,5,11,5]
    Output: true
    Explanation: The array can be partitioned as [1, 5, 5] and [11].

Example 2:
    Input: nums = [1,2,3,5]
    Output: false
    Explanation: The array cannot be partitioned into equal sum subsets.

Constraints:
    1 <= nums.length <= 200
    1 <= nums[i] <= 100
*/
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
