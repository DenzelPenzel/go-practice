/*
"""
Given an unsorted array of integers nums,
return the length of the longest consecutive elements sequence.

Follow up: Could you implement the O(n) solution?

Example 1:
  Input: nums = [100,4,200,1,3,2]
  Output: 4
  Explanation: The longest consecutive elements sequence is [1, 2, 3, 4].
    Therefore its length is 4.

Example 2:
  Input: nums = [0,3,7,2,5,8,4,6,0,1]
  Output: 9

Constraints:
  0 <= nums.length <= 104
  -231 <= nums[i] <= 231 - 1
*/

package leetcode

func longestConsecutive(nums []int) int {
	n := len(nums)
	seen := make(map[int]bool, n)
	res := 0

	for _, v := range nums {
		seen[v] = true
	}

	for _, x := range nums {
		if !seen[x-1] {
			xx := x + 1
			for seen[xx] {
				xx += 1
			}
			res = Max(res, xx-x)
		}
	}
	return res
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
