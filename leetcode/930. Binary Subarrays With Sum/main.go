package leetcode

func numSubarraysWithSum(nums []int, goal int) int {
	return atMostK(nums, goal) - atMostK(nums, goal-1)
}

func atMostK(nums []int, K int) int {
	start, end := 0, 0
	cnt_ones := 0
	res := 0
	for end < len(nums) {
		if nums[end] == 1 {
			cnt_ones += 1
		}
		end += 1
		for cnt_ones > K && start < end {
			if nums[start] == 1 {
				cnt_ones -= 1
			}
			start += 1
		}
		res += end - start
	}
	return res
}
