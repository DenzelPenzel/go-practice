/*
You are given an n x n binary matrix grid where 1 represents land and 0 represents water.

An island is a 4-directionally connected group of 1's not connected to any other 1's.

There are exactly two islands in grid.

You may change 0's to 1's to connect the two islands to form one island.

Return the smallest number of 0's you must flip to connect the two islands.

Example 1:
	Input: grid = [[0,1],[1,0]]
	Output: 1

Example 2:
	Input: grid = [[0,1,0],[0,0,0],[0,0,1]]
	Output: 2

Example 3:
	Input: grid = [[1,1,1,1,1],[1,0,0,0,1],[1,0,1,0,1],[1,0,0,0,1],[1,1,1,1,1]]
	Output: 1

Constraints:
	n == grid.length == grid[i].length
	2 <= n <= 100
	grid[i][j] is either 0 or 1.
	There are exactly two islands in grid.
*/

package leetcode

func shortestBridge(A [][]int) int {
	n, m := len(A), len(A[0])
	stack := make([][]int, 0)
	seen := make([][]bool, n)

	for i := range seen {
		seen[i] = make([]bool, m)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 1 {
				stack = append(stack, []int{i, j})
				seen[i][j] = true
				break
			}
		}
		if len(stack) > 0 {
			break
		}
	}

	for len(stack) > 0 {
		i, j := stack[len(stack)-1][0], stack[len(stack)-1][1]
		stack = stack[:len(stack)-1]
		for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j + 1}, {i, j - 1}} {
			x, y := p[0], p[1]

			if 0 <= x && x < n && 0 <= y && y < m && A[x][y] == 1 && !seen[x][y] {
				seen[x][y] = true
				stack = append(stack, []int{x, y})
			}
		}
	}

	queue := make([][]int, 0)
	for i, v := range seen {
		for j := range v {
			if v[j] {
				queue = append(queue, []int{i, j})
			}
		}
	}

	step := 0

	for len(queue) > 0 {
		s := len(queue)

		for k := 0; k < s; k++ {
			i, j := queue[0][0], queue[0][1]
			queue = queue[1:]

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j + 1}, {i, j - 1}} {
				x, y := p[0], p[1]

				if 0 <= x && x < n && 0 <= y && y < m && !seen[x][y] {
					if A[x][y] == 1 {
						return step
					}
					seen[x][y] = true
					queue = append(queue, []int{x, y})
				}
			}
		}
		step += 1
	}

	return step
}
