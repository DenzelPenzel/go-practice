/*
You are given a 0-indexed integer array nums.
You have to partition the array into one or more contiguous subarrays.

We call a partition of the array valid if each of the obtained subarrays satisfies
one of the following conditions:

The subarray consists of exactly 2 equal elements. For example, the subarray [2,2] is good.
The subarray consists of exactly 3 equal elements. For example, the subarray [4,4,4] is good.
The subarray consists of exactly 3 consecutive increasing elements,
that is, the difference between adjacent elements is 1.

For example, the subarray [3,4,5] is good, but the subarray [1,3,5] is not.

Return true if the array has at least one valid partition. Otherwise, return false.

Example 1:
	Input: nums = [4,4,4,5,6]
	Output: true
	Explanation:
		The array can be partitioned into the subarrays [4,4] and [4,5,6].
		This partition is valid, so we return true.

Example 2:
	Input: nums = [1,1,1,2]
	Output: false
	Explanation:
		There is no valid partition for this array.

Constraints:
	2 <= nums.length <= 10^5
	1 <= nums[i] <= 10^6
*/

package main

func validPartition(A []int) bool {
	n := len(A)
	dp := make([]bool, n+1)
	dp[0] = true

	for i := range A {
		eq3 := i-1 >= 0 && i-2 >= 0 && A[i] == A[i-1] && A[i] == A[i-2]
		eq2 := i-1 >= 0 && A[i] == A[i-1]
		inc := i-1 >= 0 && i-2 >= 0 && A[i] == A[i-1]+1 && A[i] == A[i-2]+2

		j := i + 1
		if eq2 && (eq3 || inc) {
			dp[j] = dp[j-3] || dp[j-2]
		} else if eq2 {
			dp[j] = dp[j-2]
		} else if eq3 || inc {
			dp[j] = dp[j-3]
		}
	}

	return dp[n]
}
