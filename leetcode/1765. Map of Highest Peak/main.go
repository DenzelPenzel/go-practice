/*
You are given an integer matrix isWater of size m x n
that represents a map of land and water cells.

If isWater[i][j] == 0, cell (i, j) is a land cell.
If isWater[i][j] == 1, cell (i, j) is a water cell.
You must assign each cell a height in a way that follows these rules:

The height of each cell must be non-negative.
If the cell is a water cell, its height must be 0.
Any two adjacent cells must have an absolute height difference of at most 1.
A cell is adjacent to another cell if the former is directly
north, east, south, or west of the latter (i.e., their sides are touching).
Find an assignment of heights such that the maximum height in the matrix is maximized.

Return an integer matrix height of size m x n where height[i][j] is cell (i, j)'s height.

If there are multiple solutions, return any of them.

Example 1:
	Input: isWater = [[0,1],[0,0]]
	Output: [[1,0],[2,1]]
	Explanation: The image shows the assigned heights of each cell.
	The blue cell is the water cell, and the green cells are the land cells.

Example 2:
	Input: isWater = [[0,0,1],[1,0,0],[0,0,0]]
	Output: [[1,1,0],[0,1,1],[1,2,2]]
	Explanation: A height of 2 is the maximum possible height of any assignment.
	Any height assignment that has a maximum height of 2 while still meeting the rules will also be accepted.


Constraints:
	m == isWater.length
	n == isWater[i].length
	1 <= m, n <= 1000
	isWater[i][j] is 0 or 1.
	There is at least one water cell.
*/

package leetcode

func highestPeak(isWater [][]int) [][]int {
	n, m := len(isWater), len(isWater[0])
	queue := make([][]int, 0)
	visited := make([][]bool, n)

	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if isWater[i][j] == 1 {
				queue = append(queue, []int{i, j})
				visited[i][j] = true
			}
		}
	}

	step := 0

	for len(queue) > 0 {
		s := len(queue)

		for i := 0; i < s; i++ {
			i, j := queue[0][0], queue[0][1]
			queue = queue[1:]
			isWater[i][j] = step

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
				x, y := p[0], p[1]

				if x >= 0 && x < n && y >= 0 && y < m && !visited[x][y] {
					visited[x][y] = true
					queue = append(queue, []int{x, y})
				}
			}
		}
		step++
	}

	return isWater
}
