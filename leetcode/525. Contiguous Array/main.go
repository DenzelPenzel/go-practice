package leetcode

func findMaxLength(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		if nums[i] == 0 {
			nums[i] = -1
		}
	}

	mapping := make(map[int]int)
	mapping[0] = -1
	prefix, target, res := 0, 0, 0

	for i, v := range nums {
		prefix += v

		if _, ok := mapping[prefix-target]; ok {
			res = Max(res, i-mapping[prefix-target])
		}

		if _, ok := mapping[prefix]; !ok {
			mapping[prefix] = i
		}
	}

	return res
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
