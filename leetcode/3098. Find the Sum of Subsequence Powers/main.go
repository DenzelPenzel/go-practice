/*
You are given an integer array nums of length n, and a positive integer k.

The power of a
subsequence
 is defined as the minimum absolute difference between any two elements in the subsequence.

Return the sum of powers of all subsequences of nums which have length equal to k.

Since the answer may be large, return it modulo 109 + 7.

Example 1:
	Input: nums = [1,2,3,4], k = 3
	Output: 4
	Explanation:
		There are 4 subsequences in nums which have length 3: [1,2,3], [1,3,4], [1,2,4], and [2,3,4].
		The sum of powers is |2 - 3| + |3 - 4| + |2 - 1| + |3 - 4| = 4.

Example 2:
	Input: nums = [2,2], k = 2
	Output: 0
	Explanation:
		The only subsequence in nums which has length 2 is [2,2].
		The sum of powers is |2 - 2| = 0.

Example 3:
	Input: nums = [4,3,-1], k = 2
	Output: 10
	Explanation:
		There are 3 subsequences in nums which have length 2: [4,3], [4,-1], and [3,-1].
		The sum of powers is |4 - 3| + |4 - (-1)| + |3 - (-1)| = 10.

Constraints:
	2 <= n == nums.length <= 50
	-108 <= nums[i] <= 108
	2 <= k <= n
*/

package main

import (
	"math"
	"sort"
)

func sumOfPowers(nums []int, k int) int {
	n := len(nums)
	mod := int(math.Pow(10, 9)) + 7
	dp := make(map[int]map[int]map[int]map[int]int, 0)
	inf := int(^uint(0) >> 1)

	var dfs func(i, last, length, diff int) int
	dfs = func(i, last, length, diff int) int {
		if length == k {
			return diff % mod
		}
		if i == n {
			return 0
		}
		if _, ok := dp[i]; !ok {
			dp[i] = make(map[int]map[int]map[int]int)
		}
		if _, ok := dp[i][length]; !ok {
			dp[i][length] = make(map[int]map[int]int)
		}
		if _, ok := dp[i][length][diff]; !ok {
			dp[i][length][diff] = make(map[int]int)
		}
		if val, ok := dp[i][length][diff][last]; ok {
			return val
		}

		var newDiff int
		if last != -1 {
			newDiff = int(math.Min(float64(diff), float64(nums[i]-nums[last])))
		} else {
			newDiff = diff
		}

		// take current element A[i] in the sequence:
		// - Update the minimum difference if necessary by comparing the absolute
		//   difference between the current element and the previous element with the current minimum difference.
		// - Proceed to the next index (i+1) and increase the length of elements taken in the sequence (count+1).
		take := dfs(i+1, i, length+1, newDiff)
		// skip current element A[i]
		notTake := dfs(i+1, last, length, diff)

		dp[i][length][diff][last] = (take%mod + notTake%mod) % mod
		return dp[i][length][diff][last]
	}

	sort.Ints(nums)
	return dfs(0, -1, 0, inf)
}
