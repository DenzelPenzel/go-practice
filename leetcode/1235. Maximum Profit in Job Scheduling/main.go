/*

We have n jobs, where every job is scheduled to be done from startTime[i] to endTime[i], obtaining a profit of profit[i].

You're given the startTime, endTime and profit arrays,
return the maximum profit you can take such that there are no two jobs in the subset with overlapping time range.

If you choose a job that ends at time X you will be able to start another job that starts at time X.


Example 1:
    Input: startTime = [1,2,3,3], endTime = [3,4,5,6], profit = [50,10,40,70]
    Output: 120
    Explanation:
        The subset chosen is the first and fourth job.
        Time range [1-3]+[3-6] , we get profit of 120 = 50 + 70.

Example 2:
    Input: startTime = [1,2,3,4,6], endTime = [3,5,10,6,9], profit = [20,20,100,70,60]
    Output: 150
    Explanation:
        The subset chosen is the first, fourth and fifth job.
        Profit obtained 150 = 20 + 70 + 60.

Example 3:
    Input: startTime = [1,1,1], endTime = [2,3,4], profit = [5,6,4]
    Output: 6

Constraints:
    1 <= startTime.length == endTime.length == profit.length <= 5 * 10^4
    1 <= startTime[i] < endTime[i] <= 10^9
    1 <= profit[i] <= 10^4
*/

package main

import "sort"

type Item struct {
	start  int
	end    int
	profit int
}

func jobScheduling(startTime []int, endTime []int, profit []int) int {
	A := make([]Item, 0)
	dp := make([]int, len(startTime))
	res := 0

	for i := 0; i < len(startTime); i++ {
		A = append(A, Item{start: startTime[i], end: endTime[i], profit: profit[i]})
	}

	sort.Slice(A, func(i, j int) bool {
		return A[i].end < A[j].end
	})

	sort.Slice(endTime, func(i, j int) bool {
		return endTime[i] < endTime[j]
	})

	search := func(target int) int {
		lo, hi := 0, len(endTime)

		for lo < hi {
			mid := lo + ((hi - lo) >> 1)

			if ends[mid] <= target {
				lo = mid + 1
			} else {
				hi = mid
			}
		}

		return lo
	}

	for i, p := range A {
		idx := search(p.start)
		prevProfit := 0

		if idx > 0 {
			prevProfit = dp[idx-1]
		}

		if i-1 >= 0 {
			dp[i] = maxI(dp[i-1], prevProfit+p.profit)
		} else {
			dp[i] = prevProfit + p.profit
		}

		res = maxI(res, dp[i])
	}

	return res
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
