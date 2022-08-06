/*
You are given an integer array nums.

You should move each element of nums into one of the two arrays A and B such
that A and B are non-empty, and average(A) == average(B).

Return true if it is possible to achieve that and false otherwise.

Note that for an array arr, average(arr) is the sum of all the elements of arr over the length of arr.

Example 1:
	Input: nums = [1,2,3,4,5,6,7,8]
	Output: true
	Explanation: We can split the array into [1,4,5,8] and [2,3,6,7], and both of them have an average of 4.5.

Example 2:
	Input: nums = [3,1]
	Output: false


Constraints:
	1 <= nums.length <= 30
	0 <= nums[i] <= 104
*/

package leetcode

func splitArraySameAverage(nums []int) bool {
	n := len(nums)
	totalSum := 0
	seen := [][][]bool{}

	for _, x := range nums {
		totalSum += x
	}

	for k := 0; k < n; k++ {
		seen = append(seen, [][]bool{})
		for i := 0; i <= totalSum; i++ {
			seen[k] = append(seen[k], []bool{})
			for j := 0; j < (n/2)+1; j++ {
				seen[k][i] = append(seen[k][i], false)
			}
		}
	}

	for subsetLen := 1; subsetLen < (n/2)+1; subsetLen++ {
		target := totalSum * subsetLen / n

		if totalSum*subsetLen%n == 0 && DFS(&nums, &seen, target, 0, subsetLen) {
			return true
		}
	}
	return false
}

func DFS(A *[]int, seen *[][][]bool, target, i, subsetLen int) bool {
	if subsetLen == 0 {
		return target == 0
	}

	if i >= len(*(A)) {
		return false
	}

	if (*seen)[target][subsetLen] {
		return true
	}

	for j := i; j < len((*A)); j++ {
		if DFS(A, seen, target-(*A)[j], j+1, subsetLen-1) {
			(*seen)[i][target][subsetLen] = true
			return true
		}
	}

	(*seen)[i][target][subsetLen] = false

	return false
}
