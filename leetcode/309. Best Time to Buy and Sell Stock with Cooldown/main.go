"""
You are given an array prices where prices[i] is the price of a given stock on the ith day.

Find the maximum profit you can achieve. 

You may complete as many transactions as you like (i.e., buy one and sell one share of the stock multiple times) with the following restrictions:

After you sell your stock, you cannot buy stock on the next day (i.e., cooldown one day).
Note: You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).

Example 1:
    Input: prices = [1,2,3,0,2]
    Output: 3
    Explanation: transactions = [buy, sell, cooldown, buy, sell]

Example 2:
    Input: prices = [1]
    Output: 0
    
Constraints:
    1 <= prices.length <= 5000
    0 <= prices[i] <= 1000
"""
package leetcode

import "math"

func maxProfit(A []int) int {
	n := len(A)
	dp := make([][2]int, n)

	dp[0][1] = -A[0] // 0 - sell, 1 - buy
	dp[0][0] = -math.MaxInt32

	for i := 1; i < n; i++ {
		dp[i][0] = Max(dp[i-1][0], dp[i-1][1]+A[i]) // sell

		if i-2 >= 0 {
			dp[i][1] = Max(dp[i-1][1], Max(dp[i-2][0]-A[i], -A[i])) // buy
		} else {
			dp[i][1] = Max(dp[i-1][1], -A[i])
		}
	}

	return Max(dp[n-1][0], 0)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
