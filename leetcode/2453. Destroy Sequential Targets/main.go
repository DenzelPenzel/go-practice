/*
You are given a 0-indexed array nums consisting of positive integers, representing targets on a number line.
You are also given an integer space.

You have a machine which can destroy targets. Seeding the machine with some nums[i] allows it to destroy
all targets with values that can be represented as nums[i] + c * space, where c is any non-negative integer.

You want to destroy the maximum number of targets in nums.

Return the minimum value of nums[i] you can seed the machine with to destroy the maximum number of targets.

Example 1:
	Input: nums = [3,7,8,1,1,5], space = 2
	Output: 1
	Explanation: If we seed the machine with nums[3], then we destroy all targets equal to 1,3,5,7,9,...
	In this case, we would destroy 5 total targets (all except for nums[2]).
	It is impossible to destroy more than 5 targets, so we return nums[3].

Example 2:
	Input: nums = [1,3,5,2,4,6], space = 2
	Output: 1
	Explanation: Seeding the machine with nums[0], or nums[3] destroys 3 targets.
	It is not possible to destroy more than 3 targets.
	Since nums[0] is the minimal integer that can destroy 3 targets, we return 1.

Example 3:
	Input: nums = [6,2,5], space = 100
	Output: 2
	Explanation: Whatever initial seed we select, we can only destroy 1 target. The minimal seed is nums[1].

Constraints:
	1 <= nums.length <= 105
	1 <= nums[i] <= 109
	1 <= space <= 109
*/

package main

import "math"

func destroyTargets(A []int, space int) int {
	mapping := make(map[int]int, 0)
	minValue := int(math.MaxInt)
	maxRemainderCnt := 0

	for _, x := range A {
		diff := x % space
		mapping[diff] += 1
	}

	for _, x := range A {
		if mapping[x%space] > maxRemainderCnt {
			maxRemainderCnt = mapping[x%space]
			minValue = x
		}

		if mapping[x%space] == maxRemainderCnt {
			minValue = min(minValue, x)
		}
	}

	return minValue
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
