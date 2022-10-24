/*
Given an array of integers arr, return the number of subarrays with an odd sum.

Since the answer can be very large, return it modulo 109 + 7.

Example 1:
    Input: arr = [1,3,5]
    Output: 4
    Explanation: All subarrays are [[1],[1,3],[1,3,5],[3],[3,5],[5]]
    All sub-arrays sum are [1,4,9,3,8,5].
    Odd sums are [1,9,3,5] so the answer is 4.

Example 2:
    Input: arr = [2,4,6]
    Output: 0
    Explanation: All subarrays are [[2],[2,4],[2,4,6],[4],[4,6],[6]]
    All sub-arrays sum are [2,6,12,4,10,6].
    All sub-arrays have even sum and the answer is 0.

Example 3:
    Input: arr = [1,2,3,4,5,6,7]
    Output: 16

Constraints:
    1 <= arr.length <= 105
    1 <= arr[i] <= 100
*/

package leetcode

func numOfSubarrays(A []int) int {
	k, mod := 2, int(1e9+7)
	prefix := 0
	mapping := make(map[int]int)
	mapping[0] = 1
	evenCnt := 0
	n := len(A)
	for _, x := range A {
		prefix += x
		d := prefix % k
		if _, ok := mapping[d]; ok {
			evenCnt += mapping[d]
		}
		mapping[d] = mapping[d] + 1
	}

	return ((n*(n+1))/2 - evenCnt) % mod
}
