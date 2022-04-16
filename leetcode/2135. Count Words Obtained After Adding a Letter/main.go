package leetcode

func wordCount(startWords []string, targetWords []string) int {
	seen := make(map[int]bool)
	res := 0

	for _, word := range startWords {
		mask := 0
		for _, c := range word {
			mask ^= 1 << (c - 'a')
		}
		seen[mask] = true
	}

	for _, word := range targetWords {
		mask := 0
		for _, c := range word {
			mask ^= 1 << (c - 'a')
		}

		for _, c := range word {
			key := mask ^ (1 << (c - 'a'))
			if seen[key] {
				res += 1
				break
			}
		}
	}

	return res

}
