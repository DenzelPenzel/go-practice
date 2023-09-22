/*
Given an integer array nums and an integer k, return the maximum length of a
subarray that sums to k. If there is not one, return 0 instead.

Example 1:

	Input: nums = [1,-1,5,-2,3], k = 3
	Output: 4
	Explanation: The subarray [1, -1, 5, -2] sums to 3 and is the longest.

Example 2:

	Input: nums = [-2,-1,2,1], k = 1
	Output: 2
	Explanation: The subarray [-1, 2] sums to 1 and is the longest.

Constraints:

	1 <= nums.length <= 2 * 105
	-104 <= nums[i] <= 104
	-109 <= k <= 109
*/
package main

import "math"

func maxSubArrayLen(nums []int, k int) int {
	current := 0
	res := 0
	mp := make(map[int]int)
	mp[0] = -1

	for i, v := range nums {
		current += v
		// Check if there exists a subarray with sum 'k' ending at index 'i'
		if prevIndex, ok := mp[current-k]; ok {
			res = int(math.Max(float64(res), float64(i-prevIndex)))
		}

		if _, ok := mp[current]; !ok {
			mp[current] = i
		}
	}

	return res
}
