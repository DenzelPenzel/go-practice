/*
You may recall that an array arr is a mountain array if and only if:

arr.length >= 3
There exists some index i (0-indexed) with 0 < i < arr.length - 1 such that:
arr[0] < arr[1] < ... < arr[i - 1] < arr[i]
arr[i] > arr[i + 1] > ... > arr[arr.length - 1]
Given an integer array arr, return the length of the longest subarray, which is a mountain.

Return 0 if there is no mountain subarray.

Example 1:
	Input: arr = [2,1,4,7,3,2,5]
	Output: 5
	Explanation: The largest mountain is [1,4,7,3,2] which has length 5.

Example 2:
	Input: arr = [2,2,2]
	Output: 0
	Explanation: There is no mountain.

Constraints:
	1 <= arr.length <= 104
	0 <= arr[i] <= 104

Follow up:
	Can you solve it using only one pass?
	Can you solve it in O(1) space?

*/

package main

func longestMountain(A []int) int {
	if len(A) <= 2 {
		return 0
	}

	res := 0
	n := len(A)
	k := 1

	for k < n-1 {
		i, j := k-1, k+1

		if j < n && A[i] < A[k] && A[k] > A[j] {
			iCnt := 0
			jCnt := 0

			for i >= 0 && A[i] < A[i+1] {
				iCnt++
				i--
			}

			for j < n && A[j-1] > A[j] {
				jCnt++
				j++
			}

			if iCnt > 0 && jCnt > 0 {
				res = max(res, iCnt+jCnt+1)
			}

			k = max(k+1, j)
		} else {
			k += 1
		}
	}

	return res
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
