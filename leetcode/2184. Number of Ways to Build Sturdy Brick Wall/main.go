package main

import (
	"fmt"
)

const MOD = 1000000007

func dfs(h int, prev int, masks []int, dp [][]int) int {
	if h == 0 {
		return 1
	}
	if dp[h][prev] == 0 {
		for _, mask := range masks {
			// there is not intersection
			if (mask & prev) == 0 {
				dp[h][prev] = (dp[h][prev] + dfs(h-1, mask, masks, dp)) % MOD
			}
		}
	}
	return dp[h][prev]
}

func buildWall(height int, width int, bricks []int) int {
	dp := make([][]int, 101)
	for i := range dp {
		dp[i] = make([]int, 1024)
	}

	var masks []int

	var perm func(w int, mask int)
	perm = func(w int, mask int) {
		if w == width {
			masks = append(masks, mask)
			return
		}
		if w != 0 {
			mask += (1 << (w - 1))
		}
		// b + w => will be as a prefix sum
		// for array [1, 1, 3] -> prefix [1, 2, 5], but total mask will be `10`
		for _, b := range bricks {
			if b+w <= width {
				perm(b+w, mask)
			}
		}
	}

	perm(0, 0)

	// return the total number of pairs
	return dfs(height, 0, masks, dp)
}

func main() {
	fmt.Println(buildWall(76, 9, []int{6, 3, 5, 1, 9}))
}
