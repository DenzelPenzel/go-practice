/*
You are given an integer array nums and an integer x.
In one operation, you can either remove the leftmost or the rightmost element from the array nums
and subtract its value from x. Note that this modifies the array for future operations.

Return the minimum number of operations to reduce x to exactly 0 if it is possible, otherwise, return -1.

Example 1:
	Input: nums = [1,1,4,2,3], x = 5
	Output: 2
	Explanation: The optimal solution is to remove the last two elements to reduce x to zero.

Example 2:
	Input: nums = [5,6,7,8,9], x = 4
	Output: -1

Example 3:
	Input: nums = [3,2,20,1,1,3], x = 10
	Output: 5
	Explanation: The optimal solution is to remove the last three elements and the first two elements (5 operations in total) to reduce x to zero.


Constraints:
	1 <= nums.length <= 10^5
	1 <= nums[i] <= 10^4
	1 <= x <= 10^9
*/

package main

import "math"

func minOperations(nums []int, target int) int {
	total := 0
	n := len(nums)
	for _, x := range nums {
		total += x
	}
	if target > total {
		return -1
	}

	A := make([]int, len(nums)*2)

	for i := 0; i < len(nums)*2; i++ {
		A[i] = nums[i%n]
	}

	ii := 0
	res := math.MaxInt32

	for i := 0; i < len(A); i++ {
		target -= A[i]

		for target < 0 {
			target += A[ii]
			ii++
		}

		if target == 0 && (ii == 0 || (ii <= n-1 && n-1 <= i)) {
			res = int(math.Min(float64(res), float64(i-ii+1)))
		}
	}

	if res == math.MaxInt32 {
		return -1
	}

	return res
}
