/*
You are given a 0-indexed 2D integer matrix grid of size 3 * 3, representing the number of stones in each cell.

The grid contains exactly 9 stones, and there can be multiple stones in a single cell.

In one move, you can move a single stone from its current cell to any other cell if the two cells share a side.

Return the minimum number of moves required to place one stone in each cell.

Example 1:
	Input: grid = [[1,1,0],[1,1,1],[1,2,1]]
	Output: 3
	Explanation:
		One possible sequence of moves to place one stone in each cell is:
		1- Move one stone from cell (2,1) to cell (2,2).
		2- Move one stone from cell (2,2) to cell (1,2).
		3- Move one stone from cell (1,2) to cell (0,2).
		In total, it takes 3 moves to place one stone in each cell of the grid.
		It can be shown that 3 is the minimum number of moves required to place one stone in each cell.

Example 2:
	Input: grid = [[1,3,0],[1,0,0],[1,0,3]]
	Output: 4
	Explanation:
		One possible sequence of moves to place one stone in each cell is:
		1- Move one stone from cell (0,1) to cell (0,2).
		2- Move one stone from cell (0,1) to cell (1,1).
		3- Move one stone from cell (2,2) to cell (1,2).
		4- Move one stone from cell (2,2) to cell (2,1).
		In total, it takes 4 moves to place one stone in each cell of the grid.
		It can be shown that 4 is the minimum number of moves required to place one stone in each cell.

Constraints:
	grid.length == grid[i].length == 3
	0 <= grid[i][j] <= 9
	Sum of grid is equal to 9.
*/

package main

import (
	"math"
)

func minimumMoves(A [][]int) int {
	var dist = func(i, j [2]int) int {
		return int(math.Abs(float64(i[0]-j[0]))) + int(math.Abs(float64(i[1]-j[1])))
	}

	n, m := len(A), len(A[0])
	coords := make([][2]int, 0)
	possible := make([][2]int, 0)
	seen := make([][]bool, n)
	for i := range seen {
		seen[i] = make([]bool, m)
	}

	var dfs func(int) int
	dfs = func(index int) int {
		if index == len(coords) {
			return 0
		}
		res := math.MaxInt32
		for _, p := range coords {
			i, j := p[0], p[1]
			if !seen[i][j] {
				seen[i][j] = true
				for _, xy := range possible {
					x, y := xy[0], xy[1]
					if A[x][y] >= 2 {
						A[x][y]--
						res = minI(res, dist([2]int{i, j}, [2]int{x, y})+dfs(index+1))
						A[x][y]++
					}
				}
				seen[i][j] = false
			}
		}
		return res
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if A[i][j] == 0 {
				coords = append(coords, [2]int{i, j})
			} else if A[i][j] >= 2 {
				possible = append(possible, [2]int{i, j})
			}
		}
	}

	return dfs(0)
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
