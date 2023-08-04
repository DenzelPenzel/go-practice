/*
You are given an array nums consisting of positive integers.

We call a subarray of an array complete if the following condition is satisfied:

The number of distinct elements in the subarray is equal to the number of distinct elements in the whole array.
Return the number of complete subarrays.

A subarray is a contiguous non-empty part of an array.

Example 1:
	Input: nums = [1,3,1,2,2]
	Output: 4
	Explanation:
		The complete subarrays are the following: [1,3,1,2], [1,3,1,2,2], [3,1,2] and [3,1,2,2].

Example 2:
	Input: nums = [5,5,5,5]
	Output: 10
	Explanation:
		The array consists only of the integer 5, so any subarray is complete.
		The number of subarrays that we can choose is 10.


Constraints:
	1 <= nums.length <= 1000
	1 <= nums[i] <= 2000
*/

func countCompleteSubarrays(nums []int) int {
	freq := make(map[int]int, 0)

	for _, v := range nums {
		if _, ok := freq[v]; !ok {
			freq[v] = 1
		} else {
			freq[v] += 1
		}
	}

	distinct := 0

	for _, v := range freq {
		if v > 0 {
			distinct++
		}
	}

	return atMost(nums, distinct) - atMost(nums, distinct-1)
}

func atMost(nums []int, k int) int {
	start, end := 0, 0
	res := 0
	unique := 0
	mp := make(map[int]int, 0)

	for end < len(nums) {
		if _, ok := mp[nums[end]]; !ok {
			unique++
		}
		mp[nums[end]] += 1
		end += 1

		for unique > k && start < end {
			mp[nums[start]] -= 1

			if mp[nums[start]] == 0 {
				unique--
				delete(mp, nums[start])
			}
			start++
		}

		res += end - start
	}

	return res
}
