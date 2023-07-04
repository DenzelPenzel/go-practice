/*
You are given a 0-indexed integer array nums. A subarray of nums is called continuous if:

Let i, i + 1, ..., j be the indices in the subarray. Then, for each pair of indices i <= i1, i2 <= j, 0 <= |nums[i1] - nums[i2]| <= 2.
Return the total number of continuous subarrays.

A subarray is a contiguous non-empty sequence of elements within an array.

Example 1:

	Input: nums = [5,4,2,4]
	Output: 8
	Explanation:
	Continuous subarray of size 1: [5], [4], [2], [4].
	Continuous subarray of size 2: [5,4], [4,2], [2,4].
	Continuous subarray of size 3: [4,2,4].
	Thereare no subarrys of size 4.
	Total continuous subarrays = 4 + 3 + 1 = 8.
	It can be shown that there are no more continuous subarrays.

Example 2:

	Input: nums = [1,2,3]
	Output: 6
	Explanation:
	Continuous subarray of size 1: [1], [2], [3].
	Continuous subarray of size 2: [1,2], [2,3].
	Continuous subarray of size 3: [1,2,3].
	Total continuous subarrays = 3 + 2 + 1 = 6.

Constraints:

	1 <= nums.length <= 105
	1 <= nums[i] <= 109
*/
package main

func continuousSubarrays(A []int) int64 {
	minDeque := make([]int, len(A))
	maxDeque := make([]int, len(A))
	res := 0
	left := 0

	// Iterate over each element in A
	for right := 0; right < len(A); right++ {
		// While the maximum deque is not empty and its last value is smaller than val, remove the last value
		for len(maxDeque) > 0 && A[maxDeque[len(maxDeque)-1]] < A[right] {
			maxDeque = maxDeque[:len(maxDeque)-1]
		}

		// Add the index of val to the max_deq
		maxDeque = append(maxDeque, right)

		// While the minimum deque is not empty and its last value is larger than val, remove the last value
		for len(minDeque) > 0 && A[minDeque[len(minDeque)-1]] > A[right] {
			minDeque = minDeque[:len(minDeque)-1]
		}

		// Add the index of val to the min_deq
		minDeque = append(minDeque, right)
		/*
			While the difference between the max and min value is greater than 2
			we remove the deque with the lowest index
			This is because we're maintaining a set of subarrays where each subarray from A[left] to A[right]
			satisfies the condition that the maximum and minimum values difference is at most 2.
			When the condition is violated, we should slide our window from the oldest value to ensure
			we're always considering new and valid subarrays. If we don't, then we're not shifting the
			'left' pointer far enough and we'll end up including some subarrays that have already been
			counted in the previous steps, leading to duplicate counts
		*/
		for len(minDeque) > 0 && len(maxDeque) > 0 && A[maxDeque[0]]-A[minDeque[0]] > 2 {
			if minDeque[0] > maxDeque[0] {
				left = maxDeque[0] + 1
				maxDeque = maxDeque[1:]
			} else {
				left = minDeque[0] + 1
				minDeque = minDeque[1:]
			}
		}
		// Count the subarrays. Every time we move the right pointer,
		// we count all valid subarrays ending at A[right], which is right - left + 1
		res += right - left + 1
	}

	return int64(res)
}
