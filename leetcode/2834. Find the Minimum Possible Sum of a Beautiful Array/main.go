/*
You are given positive integers n and target.

An array nums is beautiful if it meets the following conditions:

nums.length == n.
nums consists of pairwise distinct positive integers.
There doesn't exist two distinct indices, i and j, in the range [0, n - 1],
such that nums[i] + nums[j] == target.

Return the minimum possible sum that a beautiful array could have modulo 10^9 + 7.

Example 1:

	Input: n = 2, target = 3
	Output: 4
	Explanation:
	    We can see that nums = [1,3] is beautiful.
	    - The array nums has length n = 2.
	    - The array nums consists of pairwise distinct positive integers.
	    - There doesn't exist two distinct indices, i and j, with nums[i] + nums[j] == 3.
	    It can be proven that 4 is the minimum possible sum that a beautiful array could have.

Example 2:

	Input: n = 3, target = 3
	Output: 8
	Explanation:
	    We can see that nums = [1,3,4] is beautiful.
	    - The array nums has length n = 3.
	    - The array nums consists of pairwise distinct positive integers.
	    - There doesn't exist two distinct indices, i and j, with nums[i] + nums[j] == 3.
	    It can be proven that 8 is the minimum possible sum that a beautiful array could have.

Example 3:

	Input: n = 1, target = 1
	Output: 1
	Explanation:
	    We can see, that nums = [1] is beautiful.

Constraints:

	1 <= n <= 10^9
	1 <= target <= 10^9
*/
package main

import "math"

func minimumPossibleSum(n int, target int) int {
	// Find the sum of k numbers starting from a given number
	sumOfConsecutiveNumbers := func(start, k float64) float64 {
		return k / 2 * (2*start + (k - 1))
	}
	mod := int(math.Pow(10, 9)) + 7
	k := minI(target/2, n)

	res := sumOfConsecutiveNumbers(1.0, float64(k))

	n -= k

	if n == 1 {
		res += float64(target)
	} else if n > 1 {
		res += sumOfConsecutiveNumbers(float64(target), float64(n))
	}
	return int(res) % mod
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}
