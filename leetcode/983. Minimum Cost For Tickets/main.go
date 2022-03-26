package leetcode

import "math"

func mincostTickets(days []int, costs []int) int {
	day_set := map[int]bool{}

	for _, v := range days {
		day_set[v] = true
	}

	dp := make([]int, 366)
	for i := 1; i < 366; i++ {
		dp[i] = math.MaxInt32
	}

	dp[0] = 0
	dayMapping := map[int]int{0: 1, 1: 7, 2: 30}

	for day := 1; day < 366; day++ {
		for i, v := range costs {
			if day_set[day] {
				dp[day] = Min(dp[day], dp[Max(0, day-dayMapping[i])]+v)
			} else {
				dp[day] = dp[day-1]
			}
		}
	}

	return dp[365]
}

func Min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
