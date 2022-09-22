/*
In a string composed of 'L', 'R', and 'X' characters, like "RXXLRXRXL",
a move consists of either replacing one occurrence of "XL" with "LX",
or replacing one occurrence of "RX" with "XR". Given the starting string
start and the ending string end, return True if and only if there exists
a sequence of moves to transform one string to the other.

Example 1:
	Input: start = "RXXLRXRXL", end = "XRLXXRRLX"
	Output: true
	Explanation: We can transform start to end following these steps:
	RXXLRXRXL ->
	XRXLRXRXL ->
	XRLXRXRXL ->
	XRLXXRRXL ->
	XRLXXRRLX

Example 2:
	Input: start = "X", end = "L"
	Output: false

Constraints:
	1 <= start.length <= 104
	start.length == end.length
	Both start and end will only consist of characters in 'L', 'R', and 'X'.
*/

package leetcode

func canTransform(start string, end string) bool {
	i, j := 0, 0

	fn := func(x string) map[int]int {
		t := make(map[int]int)

		for _, v := range x {
			c := int(v) - 65
			if _, ok := t[c]; ok {
				t[c] += 1
			} else {
				t[c] = 1
			}
		}

		return t
	}

	x := fn(start)
	y := fn(end)

	if len(start) != len(end) || len(x) != len(y) {
		return false
	}

	for k, _ := range x {
		if x[k] != y[k] {
			return false
		}
	}

	for j < len(end) {

		if end[j] != 'X' {
			for i < len(start) && start[i] == 'X' {
				i += 1
			}

			if i == len(start) || start[i] != end[j] || (string(end[j]) == "L" && i < j) || (string(end[j]) == "R" && i > j) {
				return false
			}

			i += 1
		}

		j += 1

	}

	return true
}
