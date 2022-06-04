package main

import (
	"fmt"
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

func main() {
	logger := logging.New(time.RFC3339, true)

	logger.Log("info", "starting up service")
	logger.Log("warning", "no tasks found")
	logger.Log("error", "exiting: no work performed")

	fmt.Println(maximumProfit([]int{5, 4, 6, 2, 3}, []int{8, 5, 4, 3, 5}, 10))

	// fmt.Println(leetcode.GetOrder([][]int{{1, 2}, {2, 4}, {3, 2}, {4, 1}}))
}
