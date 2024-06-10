func maximumLength(A []int, k int) int {
	n := len(A)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = 1
		}
	}
	res := 1
	dp2 := make([]map[int]int, k+1)
	dp_k := make([]int, k+1)
	for i := range dp2 {
		dp2[i] = make(map[int]int)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < k+1; j++ {
			// check all possible previous transition counts
			for jj := 0; jj < j; jj++ {
				// update dp by considering previous transition counts
				dp[i][j] = max(dp[i][j], dp_k[jj]+1)
			}
			// check if the current number has been seen before at this transition count
			if idx, ok := dp2[j][A[i]]; ok {
				dp[i][j] = max(dp[i][j], idx+1)
			}
			res = max(res, dp[i][j])
		}
		for j := 0; j < k+1; j++ {
			dp2[j][A[i]] = max(dp2[j][A[i]], dp[i][j])
			dp_k[j] = max(dp_k[j], dp[i][j])
		}
	}
	return res
}

func maximumLength(A []int, k int) int {
	n := len(A)
	if n == 0 {
		return 0
	}
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	for i := 0; i < n; i++ {
		dp[i][0] = 1
	}
	res := 1
	for i := 1; i < n; i++ {
		for j := 0; j <= k; j++ {
			dp[i][j] = 1
			for p := 0; p < i; p++ {
				if A[p] == A[i] {
					dp[i][j] = max(dp[i][j], dp[p][j]+1)
				} else if j > 0 {
					dp[i][j] = max(dp[i][j], dp[p][j-1]+1)
				}
			}
			res = max(res, dp[i][j])
		}
	}

	return res
}