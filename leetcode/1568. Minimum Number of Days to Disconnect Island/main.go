/*
You are given an m x n binary grid grid where 1 represents land and 0 represents water.

An island is a maximal 4-directionally (horizontal or vertical) connected group of 1's.

The grid is said to be connected if we have exactly one island, otherwise is said disconnected.

In one day, we are allowed to change any single land cell (1) into a water cell (0).

Return the minimum number of days to disconnect the grid.

Example 1:
    Input: grid = [[0,1,1,0],[0,1,1,0],[0,0,0,0]]
        Output: 2
        Explanation: We need at least 2 days to get a disconnected grid.
        Change land grid[1][1] and grid[0][2] to water and get 2 disconnected island.

    Example 2:
        Input: grid = [[1,1]]
        Output: 2
        Explanation: Grid of full water is also disconnected ([[1,1]] -> [[0,0]]), 0 islands.

Constraints:
    m == grid.length
    n == grid[i].length
    1 <= m, n <= 30
    grid[i][j] is either 0 or 1.
*/

package main

import "math"

func minDays(A [][]int) int {
	if numIslands(A) != 1 {
		return 0
	}

	res := math.MaxInt32
	n, m := len(A), len(A[0])

	var dfs func(i, j, steps int)
	dfs = func(i, j, steps int) {
		if numIslands(A) != 1 {
			res = int(math.Min(float64(res), float64(steps)))
			return
		}

		if steps >= 2 {
			return
		}

		for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}, {i + 1, j + 1}, {i - 1, j - 1}, {i + 1, j - 1}, {i - 1, j + 1}} {
			if p[0] >= 0 && p[0] < n && p[1] >= 0 && p[1] < m && A[p[0]][p[1]] == 1 {
				A[p[0]][p[1]] = 0
				dfs(p[0], p[1], steps+1)
				A[p[0]][p[1]] = 1
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 1 {
				A[i][j] = 0
				dfs(i, j, 1)
				A[i][j] = 1
			}
		}
	}

	return res
}

func numIslands(A [][]int) int {
	n, m := len(A), len(A[0])
	seen := make([][]bool, n)
	res := 0
	for i := 0; i < n; i++ {
		seen[i] = make([]bool, m)
	}

	var dfs func(i, j int)
	dfs = func(i, j int) {
		for _, p := range [][]int{{i - 1, j}, {i + 1, j}, {i, j - 1}, {i, j + 1}} {
			if p[0] >= 0 && p[0] < n && p[1] >= 0 && p[1] < m && !seen[p[0]][p[1]] && A[p[0]][p[1]] == 1 {
				seen[p[0]][p[1]] = true
				dfs(p[0], p[1])
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 1 && !seen[i][j] {
				dfs(i, j)
				res++
			}
		}
	}

	return res
}
