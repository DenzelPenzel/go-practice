package leetcode

func targetIndices(nums []int, target int) []int {
	smaller, equal := 0, 0
	for _, v := range nums {
		if v < target {
			smaller += 1
		}
		if v == target {
			equal += 1
		}
	}

	res := make([]int, 0)

	for i := smaller; i < equal+smaller; i++ {
		res = append(res, i)
	}

	return res
}
