package main

import (
	"strconv"
	"strings"
)

/*
You are given two arrays of strings that represent two inclusive events that happened on the same day, event1 and event2, where:

event1 = [startTime1, endTime1] and
event2 = [startTime2, endTime2].

Event times are valid 24 hours format in the form of HH:MM.

A conflict happens when two events have some non-empty intersection (i.e., some moment is common to both events).

Return true if there is a conflict between two events. Otherwise, return false.

Example 1:
	Input: event1 = ["01:15","02:00"], event2 = ["02:00","03:00"]
	Output: true
	Explanation: The two events intersect at time 2:00.

Example 2:
	Input: event1 = ["01:00","02:00"], event2 = ["01:20","03:00"]
	Output: true
	Explanation: The two events intersect starting from 01:20 to 02:00.

Example 3:
	Input: event1 = ["10:00","11:00"], event2 = ["14:00","15:00"]
	Output: false
	Explanation: The two events do not intersect.

Constraints:

	evnet1.length == event2.length == 2.
	event1[i].length == event2[i].length == 5
	startTime1 <= endTime1
	startTime2 <= endTime2
	All the event times follow the HH:MM format.
*/

func haveConflict(event1 []string, event2 []string) bool {
	toMinutes := func(t string) int {
		timeParts := strings.Split(t, ":")
		start, _ := strconv.Atoi(timeParts[0])
		end, _ := strconv.Atoi(timeParts[1])
		return start*60 + end
	}

	overlap := func(x, y, u, z int) bool {
		return z >= x && y >= u
	}

	event1InMinutes := make([]int, len(event1))
	event2InMinutes := make([]int, len(event2))

	for i, x := range event1 {
		event1InMinutes[i] = toMinutes(x)
	}

	for i, e := range event2 {
		event2InMinutes[i] = toMinutes(e)
	}

	return overlap(event1InMinutes[0], event1InMinutes[1], event2InMinutes[0], event2InMinutes[1]) ||
		overlap(event2InMinutes[0], event2InMinutes[1], event1InMinutes[0], event1InMinutes[1])
}
