package leetcode

/*
Given an array of integers arr, find the sum of min(b),
where b ranges over every (contiguous) subarray of arr.

Since the answer may be large, return the answer modulo 10^9 + 7.

Example 1:

	Input: arr = [3,1,2,4]
	Output: 17
	Explanation:
	Subarrays are [3], [1], [2], [4], [3,1], [1,2], [2,4], [3,1,2], [1,2,4], [3,1,2,4].
	Minimums are 3, 1, 2, 4, 1, 1, 2, 1, 1, 1.
	Sum is 17.

Example 2:

	Input: arr = [11,81,94,43,3]
	Output: 444

Constraints:

	1 <= arr.length <= 3 * 104
	1 <= arr[i] <= 3 * 104

<=====================================SOLUTION=====================================>

If an element is smaller than x elements before and y elements after.
That element is the smallest for (x + 1) * (y + 1) subarrays.

We can model this using a monostack
Where we store positions of increasing elements.

When an element j is smaller than the element k on the top of the stack, we have (j - k) larger elements after k.
We pop element k from the stack, and now a smaller element i is on top. So, we have (k - i) larger elements before k.

	L                R

Thus, the sum of corresponding subarrays is (k - i) * A[k] * (j - k)

	0  1  2  3  4  5  6

[3, 7, 8, 5, 6, 7, 4]

	i        k        j

[7, 8, 5, 6, 7]

	^

5 is the smallest number on that range
Overall we have 9 substrings where 5 is the smallest number

[7,8,5], [8,5] [8,5,6] [5,6] [5,6,7] [7,8,5,6] [8,5,6,7] [7,8,5,6] [5] = 5 * 9 = 45 total sum
*/
func sumSubarrayMins(A []int) int {
	res, n := 0, len(A)
	mod := int(1e9 + 7)
	var st []int

	for j := 0; j < n+1; j++ {
		for len(st) > 0 && (j == n || A[st[len(st)-1]] >= A[j]) {
			k := st[len(st)-1]
			st = st[:len(st)-1]

			if len(st) > 0 {
				res += (k - st[len(st)-1]) * (j - k) * A[k]
				res %= mod
			} else {
				res += (k + 1) * (j - k) * A[k]
				res %= mod
			}
		}
		st = append(st, j)
	}

	return res

}
