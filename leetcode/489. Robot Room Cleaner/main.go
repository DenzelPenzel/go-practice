/*
You are controlling a robot that is located somewhere in a room.

The room is modeled as an m x n binary grid where 0 represents a wall and 1 represents an empty slot.

The robot starts at an unknown location in the room that is guaranteed to be empty,
 and you do not have access to the grid, but you can move the robot using the given API Robot.

You are tasked to use the robot to clean the entire room (i.e., clean every empty cell in the room).
The robot with the four given APIs can move forward, turn left, or turn right. Each turn is 90 degrees.

When the robot tries to move into a wall cell, its bumper sensor detects the obstacle,
and it stays on the current cell.

Design an algorithm to clean the entire room using the following APIs:

interface Robot {
  // returns true if next cell is open and robot moves into the cell.
  // returns false if next cell is obstacle and robot stays on the current cell.
  boolean move();

  // Robot will stay on the same cell after calling turnLeft/turnRight.
  // Each turn will be 90 degrees.
  void turnLeft();
  void turnRight();

  // Clean the current cell.
  void clean();
}
Note that the initial direction of the robot will be facing up. You can assume all four edges of the grid are all surrounded by a wall.


Custom testing:
	The input is only given to initialize the room and the robot's position internally.
	You must solve this problem "blindfolded". In other words, you must control the robot
	using only the four mentioned APIs without knowing the room layout and the initial robot's position.


Example 1:
	Input: room = [[1,1,1,1,1,0,1,1],[1,1,1,1,1,0,1,1],[1,0,1,1,1,1,1,1],[0,0,0,1,0,0,0,0],[1,1,1,1,1,1,1,1]], row = 1, col = 3
	Output: Robot cleaned all rooms.
	Explanation: All grids in the room are marked by either 0 or 1.
	0 means the cell is blocked, while 1 means the cell is accessible.
	The robot initially starts at the position of row=1, col=3.
	From the top left corner, its position is one row below and three columns right.

Example 2:
	Input: room = [[1]], row = 0, col = 0
	Output: Robot cleaned all rooms.


Constraints:
	m == room.length
	n == room[i].length
	1 <= m <= 100
	1 <= n <= 200
	room[i][j] is either 0 or 1.
	0 <= row < m
	0 <= col < n
	room[row][col] == 1
	All the empty cells can be visited from the starting position.

*/

package main

/**
 * // This is the robot's control interface.
 * // You should not implement it, or speculate about its implementation
 * type Robot struct {
 * }
 *
 * // Returns true if the cell in front is open and robot moves into the cell.
 * // Returns false if the cell in front is blocked and robot stays in the current cell.
 * func (robot *Robot) Move() bool {}
 *
 * // Robot will stay in the same cell after calling TurnLeft/TurnRight.
 * // Each turn will be 90 degrees.
 * func (robot *Robot) TurnLeft() {}
 * func (robot *Robot) TurnRight() {}
 *
 * // Clean the current cell.
 * func (robot *Robot) Clean() {}
 */

var directions = [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type Set map[[2]int]bool

func cleanRoom(robot *Robot) {
	seen := make(Set)

	goBack := func() {
		robot.TurnRight()
		robot.TurnRight()
		robot.Move()
		robot.TurnRight()
		robot.TurnRight()
	}

	var dfs func(x, y, direction int)
	dfs = func(x, y, direction int) {
		point := [2]int{x, y}
		seen[point] = true
		robot.Clean()
		for k := 0; k < 4; k++ {
			nextDirection := (direction + k) % 4
			xx := x + directions[nextDirection][0]
			yy := y + directions[nextDirection][1]

			if seen[[2]int{xx, yy}] == false && robot.Move() {
				dfs(xx, yy, nextDirection)
			}

			robot.TurnRight()
		}

		goBack()
	}

	dfs(0, 0, 0)
}
