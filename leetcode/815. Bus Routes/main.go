/*

You are given an array routes representing bus routes
where routes[i] is a bus route that the ith bus repeats forever.

For example, if routes[0] = [1, 5, 7], this means
that the 0th bus travels in the sequence 1 -> 5 -> 7 -> 1 -> 5 -> 7 -> 1 -> ... forever.
You will start at the bus stop source (You are not on any bus initially),
and you want to go to the bus stop target.

You can travel between bus stops by buses only.

Return the least number of buses you must take to travel from source to target.
Return -1 if it is not possible.

Example 1:
	Input: routes = [[1,2,7],[3,6,7]], source = 1, target = 6
	Output: 2
	Explanation: The best strategy is take the first bus to the bus stop 7, then take the second bus to the bus stop 6.

Example 2:
	Input: routes = [[7,12],[4,5,15],[6],[15,19],[9,12,13]], source = 15, target = 12
	Output: -1

Constraints:
	1 <= routes.length <= 500.
	1 <= routes[i].length <= 105
	All the values of routes[i] are unique.
	sum(routes[i].length) <= 105
	0 <= routes[i][j] < 106
	0 <= source, target < 106
*/

package main

func numBusesToDestination(routes [][]int, source int, target int) int {
	if source == target {
		return 0
	}

	// graph where the nodes are the bus routes
	graph := make(map[int][]int) // stop -> [bus...]

	for bus := 0; bus < len(routes); bus++ {
		for _, stop := range routes[bus] {
			graph[stop] = append(graph[stop], bus)
		}
	}

	queue := [][2]int{{source, 0}}

	// track visited stops
	visitedStops := make(map[int]bool)
	visitedStops[source] = true

	// track visited buses
	visitedBuses := make(map[int]bool)

	for len(queue) > 0 {
		stop, move := queue[0][0], queue[0][1]
		queue = queue[1:]

		if stop == target {
			return move
		}

		for _, bus := range graph[stop] {
			if visitedBuses[bus] {
				continue
			}

			visitedBuses[bus] = true

			for _, nxtStop := range routes[bus] {
				if visitedStops[nxtStop] {
					continue
				}
				visitedStops[nxtStop] = true
				queue = append(queue, [2]int{nxtStop, move + 1})
			}
		}
	}

	return -1
}
