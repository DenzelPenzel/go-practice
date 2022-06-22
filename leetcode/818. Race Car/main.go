/*
Your car starts at position 0 and speed +1 on an infinite number line.
Your car can go into negative positions. Your car drives automatically
according to a sequence of instructions 'A' (accelerate) and 'R' (reverse):

When you get an instruction 'A', your car does the following:
position += speed
speed *= 2
When you get an instruction 'R', your car does the following:
If your speed is positive then speed = -1
otherwise speed = 1
Your position stays the same.
For example, after commands "AAR", your car goes to positions 0 --> 1 --> 3 --> 3, and your speed goes to 1 --> 2 --> 4 --> -1.

Given a target position target, return the length of the shortest sequence of instructions to get there.

Example 1:
	Input: target = 3
	Output: 2
	Explanation:
		The shortest instruction sequence is "AA".
		Your position goes from 0 --> 1 --> 3.

Example 2:
	Input: target = 6
	Output: 5
	Explanation:
		The shortest instruction sequence is "AAARA".
		Your position goes from 0 --> 1 --> 3 --> 7 --> 7 --> 6.

Constraints:
	1 <= target <= 104
*/

package leetcode

type State struct {
	Position, Speed int
}

func Accelerate(s State) State {
	return State{
		Position: s.Position + s.Speed,
		Speed:    s.Speed * 2,
	}
}

func Reverse(s State) State {
	if s.Speed > 0 {
		return State{
			Position: s.Position,
			Speed:    -1,
		}
	}

	return State{
		Position: s.Position,
		Speed:    1,
	}
}

func racecar(target int) int {
	seen := make(map[State]bool)
	init := State{
		Position: 0,
		Speed:    1,
	}
	queue := []State{init}
	seen[init] = true
	move := 0
	actions := []func(State) State{Accelerate, Reverse}

	for len(queue) > 0 {
		size := len(queue)

		for i := 0; i < size; i++ {
			current := queue[0]
			queue = queue[1:]

			for _, action := range actions {
				next := action(current)

				if next.Position == target {
					return move + 1
				}

				if !seen[next] && next.Position >= 0 && next.Position <= 2*target {
					seen[next] = true
					queue = append(queue, next)
				}
			}

		}
		move += 1
	}

	return -1
}
