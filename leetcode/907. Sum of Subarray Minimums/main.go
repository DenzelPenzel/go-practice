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

func sumSubarrayMins(arr []int) int {
	st := make([]int, 0)
	A := make([]int, len(arr)+2)
	for i := 1; i < len(arr)+1; i++ {
		A[i] = arr[i-1]
	}

	mod := int(1e9 + 7)
	res := 0

	for j := 0; j < len(A); j++ {
		for len(st) > 0 && A[st[len(st)-1]] >= A[j] {
			k := st[len(st)-1]
			st = st[:len(st)-1]

			if len(st) > 0 {
				L := k - st[len(st)-1]
				R := j - k
				res += L * A[k] * R
			} else {
				res += A[k] * j
			}

			res %= mod
		}

		st = append(st, j)

	}

	return res
}
