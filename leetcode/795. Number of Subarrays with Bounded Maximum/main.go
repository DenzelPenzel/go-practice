package leetcode

/*
Given an integer array nums and two integers left and right,
return the number of contiguous non-empty subarrays such that
the value of the maximum array element in that subarray is in the range [left, right].

The test cases are generated so that the answer will fit in a 32-bit integer.

Example 1:
	Input: nums = [2,1,4,3], left = 2, right = 3
	Output: 3
	Explanation: There are three subarrays that meet the requirements: [2], [2, 1], [3].

Example 2:
	Input: nums = [2,9,2,5,6], left = 2, right = 8
	Output: 7

Constraints:
	1 <= nums.length <= 105
	0 <= nums[i] <= 109
	0 <= left <= right <= 109
*/

func numSubarrayBoundedMax(nums []int, left int, right int) int {
	// dp[i] = number of valid subarrays ending with i
	dp := make([]int, len(nums))
	prev := -1

	for i := 0; i < len(nums); i++ {
		if nums[i] < left {
			if i-1 >= 0 {
				dp[i] = dp[i-1]
			}
		} else if nums[i] > right {
			dp[i] = 0
			prev = i
		} else {
			dp[i] = i - prev
		}
	}

	res := 0
	for i := 0; i < len(nums); i++ {
		res += dp[i]
	}

	return res
}
