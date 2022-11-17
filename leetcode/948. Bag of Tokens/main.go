package main

import "sort"

func bagOfTokensScore(tokens []int, power int) int {
	lo, hi := 0, len(tokens)-1
	score := 0
	res := 0

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i] < tokens[j]
	})

	for lo <= hi {
		if power >= tokens[lo] {
			power -= tokens[lo]
			lo++
			score++
		} else if score > 0 {
			power += tokens[hi]
			hi--
			score--
		} else {
			break
		}

		res = Max(res, score)

	}

	return res

}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
