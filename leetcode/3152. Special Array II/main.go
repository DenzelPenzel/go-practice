/*
An array is considered special if every pair of its
adjacent elements contains two numbers with different parity.

You are given an array of integer nums and a 2D integer matrix queries,
where for queries[i] = [fromi, toi] your task is to check that
subarray nums[fromi..toi] is special or not.

Return an array of booleans answer such that answer[i] is true if nums[fromi..toi] is special.

Example 1:
	Input: nums = [3,4,1,2,6], queries = [[0,4]]
	Output: [false]
	Explanation:
		The subarray is [3,4,1,2,6]. 2 and 6 are both even.

Example 2:
	Input: nums = [4,3,1,6], queries = [[0,2],[2,3]]
	Output: [false,true]
	Explanation:
		The subarray is [4,3,1]. 3 and 1 are both odd. So the answer to this query is false.
		The subarray is [1,6]. There is only one pair: (1,6) and it contains numbers with different parity. So the answer to this query is true.


Constraints:
	1 <= nums.length <= 10^5
	1 <= nums[i] <= 10^5
	1 <= queries.length <= 10^5
	queries[i].length == 2
	0 <= queries[i][0] <= queries[i][1] <= nums.length - 1

*/

package main

func isArraySpecial(nums []int, queries [][]int) []bool {
	// modify input array, if even 1 else -1
	for i := 0; i < len(nums); i++ {
		if nums[i]%2 == 0 {
			nums[i] = 1
		} else {
			nums[i] = -1
		}
	}

	// build the prefix sum
	prefix := make([]int, 0)
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			prefix = append(prefix, 1)
		} else {
			prefix = append(prefix, 0)
		}
	}
	for i := 1; i < len(prefix); i++ {
		prefix[i] += prefix[i-1]
	}
	prefix = append([]int{0}, prefix...)
	res := make([]bool, len(queries))

	// check each query
	for i, p := range queries {
		x, y := p[0], p[1]
		if prefix[y]-prefix[x] == y-x {
			res[i] = true
		}

	}
	return res
}
