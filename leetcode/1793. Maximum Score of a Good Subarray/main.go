/*
You are given an array of integers nums (0-indexed) and an integer k.

The score of a subarray (i, j) is defined as min(nums[i], nums[i+1], ..., nums[j]) * (j - i + 1).

A good subarray is a subarray where i <= k <= j.

Return the maximum possible score of a good subarray.

Example 1:
	Input: nums = [1,4,3,7,4,5], k = 3
	Output: 15
	Explanation: The optimal subarray is (1, 5) with a score of min(4,3,7,4,5) * (5-1+1) = 3 * 5 = 15.

Example 2:
	Input: nums = [5,5,4,5,4,1,1,1], k = 0
	Output: 20
	Explanation: The optimal subarray is (0, 4) with a score of min(5,5,4,5,4) * (4-0+1) = 4 * 5 = 20.

Constraints:
	1 <= nums.length <= 105
	1 <= nums[i] <= 2 * 104
	0 <= k < nums.length
*/

package leetcode

func maximumScore(nums []int, K int) int {
	st := []int{}
	res := 0
	nums = append(nums, 0)
	nums = append([]int{0}, nums...)

	for idx := 0; idx < len(nums); idx++ {
		for len(st) > 0 && nums[st[len(st)-1]] >= nums[idx] {
			k := st[len(st)-1]
			st = st[:len(st)-1]

			if len(st) > 0 {
				i := st[len(st)-1] + 1
				j := idx - 1
				/*
					[i,....k....,j] then k is the smallest on the range between [i, k] and [k, j]
					we can calculate subarray sum in that range
				*/
				if i-1 <= K && K < j {
					res = max(res, nums[k]*(j-i+1))
				}
			} else {
				// if not elements in the stack
				if k < idx {
					res = max(res, nums[k]*(idx-1))
				}
			}
		}
		// append current element in the stack
		st = append(st, idx)
	}

	return res

}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
