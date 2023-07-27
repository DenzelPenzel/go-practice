/*
(This problem is an interactive problem.)

You may recall that an array arr is a mountain array if and only if:

arr.length >= 3
There exists some i with 0 < i < arr.length - 1 such that:
arr[0] < arr[1] < ... < arr[i - 1] < arr[i]
arr[i] > arr[i + 1] > ... > arr[arr.length - 1]
Given a mountain array mountainArr, return the minimum index such that mountainArr.get(index) == target.
If such an index does not exist, return -1.

You cannot access the mountain array directly.
You may only access the array using a MountainArray interface:

MountainArray.get(k) returns the element of the array at index k (0-indexed).
MountainArray.length() returns the length of the array.
Submissions making more than 100 calls to MountainArray.get will be judged Wrong Answer.
Also, any solutions that attempt to circumvent the judge will result in disqualification.

Example 1:
	Input: array = [1,2,3,4,5,3,1], target = 3
	Output: 2
	Explanation: 3 exists in the array, at index=2 and index=5. Return the minimum index, which is 2.

Example 2:
	Input: array = [0,1,2,4,2,1], target = 3
	Output: -1
	Explanation: 3 does not exist in the array, so we return -1.

Constraints:
	3 <= mountain_arr.length() <= 104
	0 <= target <= 109
	0 <= mountain_arr.get(index) <= 109
*/

package main

/**
 * // This is the MountainArray's API interface.
 * // You should not implement it, or speculate about its implementation
 * type MountainArray struct {
 * }
 *
 * func (this *MountainArray) get(index int) int {}
 * func (this *MountainArray) length() int {}
 */

func findInMountainArray(target int, mountainArr *MountainArray) int {
	n := mountainArr.length()
	lo, hi := 0, n

	for lo < hi {
		mid := lo + ((hi - lo) >> 1)

		if mid+1 < n && mountainArr.get(mid) <= mountainArr.get(mid+1) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	var search func(lo, hi int, reversed bool) int
	search = func(lo, hi int, reversed bool) int {
		for lo < hi {
			mid := lo + ((hi - lo) >> 1)
			x := mountainArr.get(mid)

			if x == target {
				return mid
			}

			if !reversed {
				if x < target {
					lo = mid + 1
				} else {
					hi = mid
				}
			} else {
				if x < target {
					hi = mid
				} else {
					lo = mid + 1
				}
			}
		}

		return lo
	}

	leftmostVal := mountainArr.get(0)
	pivotVal := mountainArr.get(lo)
	pivotIdx := lo

	if pivotVal == target {
		return pivotIdx
	}

	if leftmostVal <= target && target <= pivotVal {
		lo := search(0, pivotIdx, false)
		if lo != n && mountainArr.get(lo) == target {
			return lo
		}
	}

	rightmostVal := mountainArr.get(n - 1)

	if pivotVal >= target && target >= rightmostVal {
		lo := search(pivotIdx+1, n, true)
		if lo != n && mountainArr.get(lo) == target {
			return lo
		}
	}

	return -1
}
