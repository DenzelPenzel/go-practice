/*
An array is considered special if every pair of its
adjacent elements contains two numbers with different parity.

You are given an array of integer a and a 2D integer matrix queries,
where for queries[i] = [fromi, toi] your task is to check that
subarray a[fromi..toi] is special or not.

Return an array of booleans answer such that answer[i] is true if a[fromi..toi] is special.

Example 1:
	Input: a = [3,4,1,2,6], queries = [[0,4]]
	Output: [false]
	Explanation:
		The subarray is [3,4,1,2,6]. 2 and 6 are both even.

Example 2:
	Input: a = [4,3,1,6], queries = [[0,2],[2,3]]
	Output: [false,true]
	Explanation:
		The subarray is [4,3,1]. 3 and 1 are both odd. So the answer to this query is false.
		The subarray is [1,6]. There is only one pair: (1,6) and it contains numbers with different parity. So the answer to this query is true.


Constraints:
	1 <= a.length <= 10^5
	1 <= a[i] <= 10^5
	1 <= queries.length <= 10^5
	queries[i].length == 2
	0 <= queries[i][0] <= queries[i][1] <= a.length - 1

*/

package main

func isArraySpecial(a []int, queries [][]int) []bool {
	// modify input array, if even 1 else -1
	for i := 0; i < len(a); i++ {
		if a[i]%2 == 0 {
			a[i] = 1
		} else {
			a[i] = -1
		}
	}

	// build the prefix sum
	prefix := make([]int, 0)
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
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

func maximumLength(A []int, k int) int {
	n := len(A)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	for i := 0; i < n; i++ {
		dp[i][0] = 1
	}
	res := 1
	for i := 1; i < n; i++ {
		for j := 0; j <= k; j++ {
			dp[i][j] = 1
			for p := 0; p < i; p++ {
				if A[p] == A[i] {
					dp[i][j] = max(dp[i][j], dp[p][j]+1)
				} else if j > 0 {
					dp[i][j] = max(dp[i][j], dp[p][j-1]+1)
				}
			}
			res = max(res, dp[i][j])
		}
	}

	return res
}

func maximumLengthFunc(A []int, k int) int {
	n := len(A)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = 1
		}
	}

	res := 1
	dp2 := make([]map[int]int, k+1)
	dp_k := make([]int, k+1)

	for i := range dp2 {
		dp2[i] = make(map[int]int)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < k+1; j++ {
			for jj := 0; jj < j; jj++ {
				dp[i][j] = max(dp[i][j], 1+dp_k[jj])
			}
			if idx, ok := dp2[j][A[i]]; ok {
				dp[i][j] = max(dp[i][j], 1+idx)
			}
			res = max(res, dp[i][j])
		}
		for j := 0; j < k+1; j++ {
			dp_k[j] = max(dp_k[j], dp[i][j])
			dp2[j][A[i]] = max(dp2[j][A[i]], dp[i][j])
		}
	}
	return res
}
