/*
You are given an array nums of non-negative integers and an integer k.

An array is called special if the bitwise OR of all of its elements is at least k.

Return the length of the shortest special non-empty
subarray
 of nums, or return -1 if no special subarray exists.

Example 1:
	Input: nums = [1,2,3], k = 2
	Output: 1
	Explanation:
		The subarray [3] has OR value of 3. Hence, we return 1.

Example 2:
	Input: nums = [2,1,8], k = 10
	Output: 3
	Explanation:
		The subarray [2,1,8] has OR value of 11. Hence, we return 3.

Example 3:
	Input: nums = [1,2], k = 0
	Output: 1
	Explanation:
		The subarray [1] has OR value of 1. Hence, we return 1.


Constraints:
	1 <= nums.length <= 2 * 10^5
	0 <= nums[i] <= 10^9
	0 <= k <= 10^9

*/

package main

func minimumSubarrayLength(A []int, k int) int {
	bits := make([]int, 32)
	lo, hi := 0, 0
	x := A[0]
	inf := int(^uint(0) >> 1)
	res := inf
	for hi < len(A) {
		x |= A[hi]
		add(&bits, A[hi])
		for lo <= hi && x >= k {
			res = min(res, hi-lo+1)
			remove(&bits, A[lo])
			x = getNum(bits)
			lo++
		}
		hi++
	}

	if res == inf {
		return -1
	}

	return res
}

func add(bits *[]int, x int) {
	for i := 0; i < 32; i++ {
		if x&(1<<i) >= 1 {
			(*bits)[i] += 1
		}
	}
}

func remove(bits *[]int, x int) {
	for i := 0; i < 32; i++ {
		if x&(1<<i) >= 1 {
			(*bits)[i] -= 1
		}
	}
}

func getNum(bits []int) int {
	res := 0
	for i, b := range bits {
		if b > 0 {
			res |= (1 << i)
		}
	}
	return res
}
