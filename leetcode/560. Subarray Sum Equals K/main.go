/*
Given an array of integers nums and an integer k,
return the total number of subarrays whose sum equals to k.

Example 1:
	Input: nums = [1,1,1], k = 2
	Output: 2

Example 2:
	Input: nums = [1,2,3], k = 3
	Output: 2

Constraints:
	1 <= nums.length <= 2 * 104
	-1000 <= nums[i] <= 1000
	-107 <= k <= 107
*/

package leetcode

func subarraySum(nums []int, target int) int {
	mapping := make(map[int]int)
	mapping[0] = 1
	prefix := 0
	res := 0

	for _, v := range nums {
		prefix += v
		if _, ok := mapping[prefix-target]; ok {
			res += mapping[prefix-target]
		}
		mapping[prefix]++
	}

	return res
}
