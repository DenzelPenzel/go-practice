package leetcode

import "strings"

func countTexts(pressedKeys string) int {
	n := len(pressedKeys)
	dp := make([]int, n)
	mod := int(1e9 + 7)

	for i := 0; i < n; i++ {
		dp[i] = 1
	}

	for i := n - 2; i >= 0; i-- {
		res := 0
		x := pressedKeys[i:]

		for _, key := range []string{"222", "22", "2", "333", "33", "3", "444", "44", "4", "555", "55", "5", "666", "66", "6", "7777", "777", "77", "7", "888", "88", "8", "9999", "999", "99", "9"} {
			if strings.HasPrefix(x, key) {
				j := i + len(key)
				if j < len(dp) {
					res += dp[j]
				} else {
					res += 1
				}
				res %= mod
				if len(key) == 1 {
					break
				}
			}
		}
		dp[i] = res
	}

	return dp[0]
}
