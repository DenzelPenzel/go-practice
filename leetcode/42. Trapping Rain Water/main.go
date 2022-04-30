/*

Given n non-negative integers representing an elevation map
where the width of each bar is 1,
compute how much water it can trap after raining.

Example 1:
  Input: height = [0,1,0,2,1,0,1,3,2,1,2,1]
  Output: 6
  Explanation: The above elevation map (black section) is represented by array [0,1,0,2,1,0,1,3,2,1,2,1]. In this case, 6 units of rain water (blue section) are being trapped.

Example 2:
  Input: height = [4,2,0,3,2,5]
  Output: 9

Constraints:
  n == height.length
  0 <= n <= 3 * 104
  0 <= height[i] <= 105
*/

package leetcode

func trap(A []int) int {
	st := make([]int, 0)
	res := 0

	for i := 0; i < len(A); i++ {
		for len(st) > 0 && A[st[len(st)-1]] < A[i] {
			j := st[len(st)-1]
			st = st[:len(st)-1]

			if len(st) > 0 {
				w := i - st[len(st)-1] - 1
				h := min(A[st[len(st)-1]], A[i]) - A[j]
				res += w * h
			}
		}
		st = append(st, i)
	}
	return res
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
