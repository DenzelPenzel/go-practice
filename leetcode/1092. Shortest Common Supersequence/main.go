package leetcode

func shortestCommonSupersequence(str1 string, str2 string) string {
	n, m := len(str1), len(str2)
	dp := [][]int{}

	for i := 0; i <= n; i++ {
		dp = append(dp, []int{})
		for j := 0; j <= m; j++ {
			dp[i] = append(dp[i], 0)
		}
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if str1[i-1] == str2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = Max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	i, j := n, m
	lcs := ""

	for i > 0 && j > 0 {
		if str1[i-1] == str2[j-1] {
			lcs += string(str1[i-1])
			i, j = i-1, j-1
		} else if dp[i-1][j] > dp[i][j-1] {
			lcs += string(str1[i-1])
			i--
		} else {
			lcs += string(str2[j-1])
			j--
		}
	}

	for i > 0 {
		lcs += string(str1[i-1])
		i--
	}

	for j > 0 {
		lcs += string(str2[j-1])
		j--
	}

	return Reverse(lcs)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
