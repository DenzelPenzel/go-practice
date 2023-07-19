/*
Given an integer array nums, find a subarray that has the largest product, and return the product.

The test cases are generated so that the answer will fit in a 32-bit integer.

Example 1:
    Input: nums = [2,3,-2,4]
    Output: 6
    Explanation: [2,3] has the largest product 6.

Example 2:
    Input: nums = [-2,0,-1]
    Output: 0
    Explanation: The result cannot be 2, because [-2,-1] is not a subarray.

Constraints:
    1 <= nums.length <= 2 * 104
    -10 <= nums[i] <= 10
    The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.
*/

package main

func maxProduct(nums []int) int {
	minVal, maxVal := nums[0], nums[0]
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i]*minVal > nums[i]*maxVal {
			minVal, maxVal = maxVal, minVal
		}
		minVal = minI(minVal*nums[i], nums[i])
		maxVal = maxI(maxVal*nums[i], nums[i])
		res = maxI(res, maxVal)
	}

	return res
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
