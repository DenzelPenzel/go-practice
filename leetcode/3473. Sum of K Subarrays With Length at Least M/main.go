/*

You are given an integer array nums and two integers, k and m.

Return the maximum sum of k non-overlapping subarrays of nums, where each subarray has a length of at least m.

Example 1:
	Input: nums = [1,2,-1,3,3,4], k = 2, m = 2
	Output: 13
	Explanation:
		The optimal choice is:

		Subarray nums[3..5] with sum 3 + 3 + 4 = 10 (length is 3 >= m).
		Subarray nums[0..1] with sum 1 + 2 = 3 (length is 2 >= m).
		The total sum is 10 + 3 = 13.

Example 2:
	Input: nums = [-10,3,-1,-2], k = 4, m = 1
	Output: -10
	Explanation:
		The optimal choice is choosing each element as a subarray. The output is (-10) + 3 + (-1) + (-2) = -10.


Constraints:
	1 <= nums.length <= 2000
	-104 <= nums[i] <= 104
	1 <= k <= floor(nums.length / m)
	1 <= m <= 3
*/

package main

func maxSum(nums []int, k int, m int) int {
	n := len(nums)
	prefix := make([]int, n+1)

	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + nums[i]
	}

	memo := make(map[[3]int]int)

	var dfs func(i, k, extend int) int
	dfs = func(i, k, extend int) int {
		if k == 0 && extend == -1 {
			return 0
		}
		if i >= n {
			if k == 0 {
				return 0
			}
			return -1e9
		}
		if i+(k*m) > n {
			return -1e9
		}

		key := [3]int{i, k, extend}
		if val, found := memo[key]; found {
			return val
		}

		res := dfs(i+1, k, -1)

		if extend == 1 {
			res = max(res, nums[i]+dfs(i+1, k, extend))
		}

		if k > 0 && i+m <= n {
			res = max(res, prefix[i+m]-prefix[i]+dfs(i+m, k-1, 1))
		}

		memo[key] = res
		return res
	}

	result := dfs(0, k, -1)
	return result
}
