/*
Given an integer array nums and an integer k,
return the k most frequent elements. You may return the answer in any order.

Example 1:
	Input: nums = [1,1,1,2,2,3], k = 2
	Output: [1,2]

Example 2:
	Input: nums = [1], k = 1
	Output: [1]

Constraints:
	1 <= nums.length <= 10^5
	-104 <= nums[i] <= 10^4
	k is in the range [1, the number of unique elements in the array].
	It is guaranteed that the answer is unique.


Follow up:
	Your algorithm's time complexity must be better than O(n log n),
	where n is the array's size.
*/

package main

func topKFrequent(nums []int, k int) []int {
	count := make(map[int]int, 0)
	for _, v := range nums {
		count[v] += 1
	}

	bucket := make([][]int, len(nums)+1)

	for k, v := range count {
		bucket[v] = append(bucket[v], k)
	}

	res := make([]int, 0)

	for i := len(nums); i >= 1; i-- {
		for len(bucket[i]) > 0 && len(res) < k {
			val := bucket[i][len(bucket[i])-1]
			res = append(res, val)
			bucket[i] = bucket[i][:len(bucket[i])-1]
		}
		if len(res) == k {
			break
		}
	}

	return res
}
