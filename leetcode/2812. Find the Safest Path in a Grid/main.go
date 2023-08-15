/*
2812. Find the Safest Path in a Grid (Medium)

    You are given a 0-indexed 2D matrix grid of size n x n, where (r, c)
    represents:
    * A cell containing a thief if grid[r][c] = 1
    * An empty cell if grid[r][c] = 0
    You are initially positioned at cell (0, 0). In one move, you can move to
    any adjacent cell in the grid, including cells containing thieves. The
    safeness factor of a path on the grid is defined as the minimum manhattan
    distance from any cell in the path to any thief in the grid. Return the
    maximum safeness factor of all paths leading to cell (n - 1, n - 1). An
    adjacent cell of cell (r, c), is one of the cells (r, c + 1), (r, c - 1),
    (r + 1, c) and (r - 1, c) if it exists. The Manhattan distance between two
    cells (a, b) and (x, y) is equal to |a - x| + |b - y|, where |val| denotes
    the absolute value of val.

    Example 1:
    Input: grid = [[1,0,0],[0,0,0],[0,0,1]]
    Output: 0
    Explanation: All paths from (0, 0) to (n - 1, n - 1) go through the thieves
                 in cells (0, 0) and (n - 1, n - 1).

    Example 2:
    Input: grid = [[0,0,1],[0,0,0],[0,0,0]]
    Output: 2
    Explanation: The path depicted in the picture above has a safeness factor
                 of 2 since:
                 - The closest cell of the path to the thief at cell (0, 2) is
                   cell (0, 0). The distance between them is
                   | 0 - 0 | + | 0 - 2 | = 2.
                 It can be shown that there are no other paths with a higher
                 safeness factor.

    Example 3:
    Input: grid = [[0,0,0,1],[0,0,0,0],[0,0,0,0],[1,0,0,0]]
    Output: 2
    Explanation: The path depicted in the picture above has a safeness factor
                 of 2 since:
                 - The closest cell of the path to the thief at cell (0, 3) is
                   cell (1, 2). The distance between them is
                   | 0 - 1 | + | 3 - 2 | = 2.
                 - The closest cell of the path to the thief at cell (3, 0) is
                   cell (3, 2). The distance between them is
                   | 3 - 3 | + | 0 - 2 | = 2.
                 It can be shown that there are no other paths with a higher
                 safeness factor.

    Constraints:
    * 1 <= grid.length == n <= 400
    * grid[i].length == n
    * grid[i][j] is either 0 or 1.
    * There is at least one thief in the grid.
*/

package main

import (
	"container/heap"
	"fmt"
	"math"
)

func maximumSafenessFactor(grid [][]int) int {
	n, m := len(grid), len(grid[0])
	queue := make([][]int, 0)
	dist := make([][]int, n)

	for i := 0; i < n; i++ {
		dist[i] = make([]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = math.MaxInt64
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 1 {
				queue = append(queue, []int{i, j})
				dist[i][j] = 0
			}
		}
	}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		i, j := p[0], p[1]

		for _, p2 := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
			x, y := p2[0], p2[1]
			if x >= 0 && x < n && y >= 0 && y < m {
				if dist[x][y] > dist[i][j]+1 {
					dist[x][y] = dist[i][j] + 1
					queue = append(queue, []int{x, y})
				}
			}
		}
	}

	fmt.Println(dist)

	maxQueue := make(PQ, 0)
	heap.Push(&maxQueue, &Node{i: 0, j: 0, minDist: dist[0][0], currentDist: dist[0][0]})

	vis := make([][]bool, n)
	for i := 0; i < n; i++ {
		vis[i] = make([]bool, m)
	}

	for len(maxQueue) > 0 {
		p := heap.Pop(&maxQueue).(*Node)
		i, j := p.i, p.j

		if i == n-1 && j == m-1 {
			return p.minDist
		}

		for _, p2 := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
			x, y := p2[0], p2[1]
			if x >= 0 && x < n && y >= 0 && y < m && !vis[x][y] {
				vis[x][y] = true
				heap.Push(&maxQueue, &Node{i: x, j: y, minDist: minI(p.minDist, dist[x][y]), currentDist: dist[x][y]})
			}
		}
	}

	return -1
}

// PriorityQueue imp
type Node struct {
	minDist     int
	currentDist int
	i, j        int
}

type PQ []*Node

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Less(i, j int) bool {
	return pq[i].currentDist > pq[j].currentDist
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PQ) Push(p interface{}) {
	*pq = append(*pq, p.(*Node))
}

func (pq *PQ) Pop() interface{} {
	val := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return val
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
