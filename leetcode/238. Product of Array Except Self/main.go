/*
238. Product of Array Except Self (Medium)

Given an array nums of n integers where n > 1,
return an array output such that output[i]
is equal to the product of all the elements of nums except nums[i].

The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.

You must write an algorithm that runs in O(n) time and without using the division operation.

Example 1:
    Input:  [1,2,3,4]
    Output: [24,12,8,6]

Example 2:
    Input: nums = [-1,1,0,-3,3]
    Output: [0,0,9,0,0]

Constraint:
    2 <= nums.length <= 105
    -30 <= nums[i] <= 30
    The product of any prefix or suffix of nums is guaranteed to fit in a 32-bit integer.

Note: Please solve it without division and in O(n).

Follow up:
    Can you solve the problem in O(1) extra space complexity?
    (The output array does not count as extra space for space complexity analysis.)
*/

package leetcode

func productExceptSelf(nums []int) []int {
	n := len(nums)
	left := make([]int, n)
	res := make([]int, n)
	left[0] = nums[0]
	for i := 1; i < n; i++ {
		left[i] = left[i-1] * nums[i]
	}
	right := 1
	for i := n - 1; i >= 0; i-- {
		if i == 0 {
			res[i] = right
		} else {
			res[i] = left[i-1] * right
			right *= nums[i]
		}
	}
	return res
}
