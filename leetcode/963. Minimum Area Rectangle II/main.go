/*

You are given an array of points in the X-Y plane points where points[i] = [xi, yi].

Return the minimum area of any rectangle formed from these points,
with sides not necessarily parallel to the X and Y axes.

If there is not any such rectangle, return 0.

Answers within 10-5 of the actual answer will be accepted.

Example 1:
	Input: points = [[1,2],[2,1],[1,0],[0,1]]
	Output: 2.00000
	Explanation: The minimum area rectangle occurs at [1,2],[2,1],[1,0],[0,1], with an area of 2.

Example 2:
	Input: points = [[0,1],[2,1],[1,1],[1,0],[2,0]]
	Output: 1.00000
	Explanation: The minimum area rectangle occurs at [1,0],[1,1],[2,1],[2,0], with an area of 1.

Example 3:
	Input: points = [[0,3],[1,2],[3,1],[1,3],[2,1]]
	Output: 0
	Explanation: There is no possible rectangle to form from these points.

Constraints:
	1 <= points.length <= 50
	points[i].length == 2
	0 <= xi, yi <= 4 * 104
	All the given points are unique.

*/
package leetcode

import "math"

func minAreaFreeRect(points [][]int) float64 {
	getDist := func(p1, p2 []int) int {
		dx := p1[0] - p2[0]
		dy := p1[1] - p2[1]
		return dx*dx + dy*dy
	}

	mapping := make(map[[3]int][][2]int)
	n := len(points)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p1 := points[i]
			p2 := points[j]

			dist := getDist(p1, p2)

			var key = [3]int{dist, p1[0] + p2[0], p1[1] + p2[1]}

			v := mapping[key]
			v = append(v, [2]int{i, j})

			mapping[key] = v
		}
	}

	res := 0

	for _, v := range mapping {

		for i := 0; i < len(v); i++ {
			for j := i + 1; j < len(v); j++ {
				p1, p2, p3 := points[v[i][0]], points[v[j][0]], points[v[j][1]]

				area := getDist(p1, p2) * getDist(p1, p3)

				if res == 0 || res > area {
					res = area
				}
			}
		}

	}

	return math.Sqrt(float64(res))
}
