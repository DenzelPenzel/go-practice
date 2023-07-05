/*
Given an array of integers nums and an integer limit,
return the size of the longest non-empty subarray such that
the absolute difference between any two elements of this subarray is less than or equal to limit.

Example 1:

	Input: nums = [8,2,4,7], limit = 4
	Output: 2
	Explanation: All subarrays are:
	[8] with maximum absolute diff |8-8| = 0 <= 4.
	[8,2] with maximum absolute diff |8-2| = 6 > 4.
	[8,2,4] with maximum absolute diff |8-2| = 6 > 4.
	[8,2,4,7] with maximum absolute diff |8-2| = 6 > 4.
	[2] with maximum absolute diff |2-2| = 0 <= 4.
	[2,4] with maximum absolute diff |2-4| = 2 <= 4.
	[2,4,7] with maximum absolute diff |2-7| = 5 > 4.
	[4] with maximum absolute diff |4-4| = 0 <= 4.
	[4,7] with maximum absolute diff |4-7| = 3 <= 4.
	[7] with maximum absolute diff |7-7| = 0 <= 4.
	Therefore, the size of the longest subarray is 2.

Example 2:

	Input: nums = [10,1,2,4,7,2], limit = 5
	Output: 4
	Explanation: The subarray [2,4,7,2] is the longest since the maximum absolute diff is |2-7| = 5 <= 5.

Example 3:

	Input: nums = [4,2,2,2,4,4,2,2], limit = 0
	Output: 3

Constraints:

	1 <= nums.length <= 105
	1 <= nums[i] <= 109
	0 <= limit <= 109
*/
package main

func longestSubarray(A []int, limit int) int {
	minDeque := make([]int, len(A))
	maxDeque := make([]int, len(A))
	res := 0
	left := 0

	for right := 0; right < len(A); right++ {
		for len(maxDeque) > 0 && A[maxDeque[len(maxDeque)-1]] < A[right] {
			maxDeque = maxDeque[:len(maxDeque)-1]
		}

		maxDeque = append(maxDeque, right)

		for len(minDeque) > 0 && A[minDeque[len(minDeque)-1]] > A[right] {
			minDeque = minDeque[:len(minDeque)-1]
		}

		minDeque = append(minDeque, right)

		for len(minDeque) > 0 && len(maxDeque) > 0 && A[maxDeque[0]]-A[minDeque[0]] > limit {
			if minDeque[0] > maxDeque[0] {
				left = maxDeque[0] + 1
				maxDeque = maxDeque[1:]
			} else {
				left = minDeque[0] + 1
				minDeque = minDeque[1:]
			}
		}

		res = maxI(res, right-left+1)
	}

	return res
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
