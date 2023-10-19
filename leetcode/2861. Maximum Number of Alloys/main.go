/*


 */

package main

import "math"

func maxNumberOfAlloys(n int, k int, budget int, composition [][]int, stock []int, costs []int) int {
	// Define a function to calculate the minimum cost for a given number of alloys.
	getMinCost := func(alloys int) int {
		res := math.MaxInt32

		for i := 0; i < k; i++ {
			cost := 0
			for j := 0; j < n; j++ {
				cost += maxI(composition[i][j]*alloys-stock[j], 0) * costs[j]
			}

			res = minI(res, cost)
		}

		return res
	}

	lo, hi := 0, int(math.Pow(10, 9))

	for lo < hi {
		mid := lo + ((hi - lo) >> 1)
		minCost := getMinCost(mid)
		if minCost <= budget {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo - 1

}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
