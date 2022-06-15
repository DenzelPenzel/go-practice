package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/denisschmidt/go-leetcode/logging"
)

func maximumProfit(present []int, future []int, budget int) int {
	A := make([][]int, 0)
	for i := 0; i < len(future); i++ {
		if future[i]-present[i] > 0 {
			A = append(A, []int{present[i], future[i] - present[i]})
		}
	}

	if len(A) == 0 {
		return 0
	}

	n := len(A)
	dp := [][]int{}

	for i := 0; i < n; i++ {
		dp = append(dp, []int{})
		for j := 0; j <= budget; j++ {
			dp[i] = append(dp[i], 0)
		}
	}

	for w := 0; w <= budget; w++ {
		if A[0][0] <= w {
			dp[0][w] = A[0][1]
		}
	}

	for i := 1; i < n; i++ {
		for w := 0; w <= budget; w++ {
			if A[i][0] <= w {
				dp[i][w] = Max(dp[i-1][w], dp[i-1][w-A[i][0]]+A[i][1])
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	return dp[len(dp)-1][budget]
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type Shop struct {
	dist  int
	price int
	row   int
	col   int
}

func highestRankedKItems(grid [][]int, pricing []int, start []int, k int) [][]int {
	n, m := len(grid), len(grid[0])
	seen := make([][]bool, n)
	queue := make([][]int, 0)
	queue = append(queue, []int{start[0], start[1]})
	items := make([]Shop, 0)
	dist := 0

	for i := range seen {
		seen[i] = make([]bool, m)
	}

	for len(queue) > 0 {
		size := len(queue)

		for k := 0; k < size; k++ {
			i, j := queue[0][0], queue[0][1]
			queue = queue[1:]

			if pricing[0] <= grid[i][j] && grid[i][j] <= pricing[1] {
				items = append(items, Shop{dist, grid[i][j], i, j})
			}

			for _, p := range [][]int{{i + 1, j}, {i - 1, j}, {i, j - 1}, {i, j + 1}} {
				x, y := p[0], p[1]

				if x >= 0 && x < n && y >= 0 && y < m {
					if grid[x][y] >= 1 {
						queue = append(queue, []int{x, y})
						seen[x][y] = true
					}
				}
			}
		}
		dist++
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].dist != items[j].dist {
			return items[i].dist < items[j].dist
		}
		if items[i].price != items[j].price {
			return items[i].price < items[j].price
		}
		if items[i].row != items[j].row {
			return items[i].row < items[j].row
		}
		return items[i].col < items[j].col
	})

	res := make([][]int, 0, k)
	for i := 0; i < len(items) && i < k; i++ {
		res = append(res, []int{items[i].row, items[i].col})
	}
	return res
}

func main() {
	logger := logging.New(time.RFC3339, true)

	logger.Log("info", "starting up service")
	logger.Log("warning", "no tasks found")
	logger.Log("error", "exiting: no work performed")

	fmt.Println(maximumProfit([]int{5, 4, 6, 2, 3}, []int{8, 5, 4, 3, 5}, 10))

	// fmt.Println(leetcode.GetOrder([][]int{{1, 2}, {2, 4}, {3, 2}, {4, 1}}))
}

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
		size := len(queue)

		for k := 0; k < size; k++ {
			i, j := queue[0][0], queue[0][1]
			queue := queue[1:]

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
