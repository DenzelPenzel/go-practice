/*
Alice is texting Bob using her phone.

The mapping of digits to letters is shown in the figure below.

In order to add a letter, Alice has to press the key of the corresponding digit i times,
where i is the position of the letter in the key.

For example, to add the letter 's', Alice has to press '7' four times.
Similarly, to add the letter 'k', Alice has to press '5' twice.

Note that the digits '0' and '1' do not map to any letters, so Alice does not use them.
However, due to an error in transmission, Bob did not receive Alice's text message
but received a string of pressed keys instead.

For example, when Alice sent the message "bob", Bob received the string "2266622".
Given a string pressedKeys representing the string received by Bob,

return the total number of possible text messages Alice could have sent.

Since the answer may be very large, return it modulo 109 + 7.

Example 1:
    Input: pressedKeys = "22233"
    Output: 8
    Explanation:
    The possible text messages Alice could have sent are:
    "aaadd", "abdd", "badd", "cdd", "aaae", "abe", "bae", and "ce".
    Since there are 8 possible messages, we return 8.

Example 2:
    Input: pressedKeys = "222222222222222222222222222222222222"
    Output: 82876089
    Explanation:
    There are 2082876103 possible text messages Alice could have sent.
    Since we need to return the answer modulo 109 + 7, we return 2082876103 % (109 + 7) = 82876089.

Constraints:
    1 <= pressedKeys.length <= 105
    pressedKeys only consists of digits from '2' - '9'.
*/
package leetcode

import "strings"

func countTexts(pressedKeys string) int {
	n := len(pressedKeys)
	dp := make([]int, n)
	mod := int(1e9 + 7)

	for i := 0; i < n; i++ {
		dp[i] = 1
	}

	for i := n - 2; i >= 0; i-- {
		res := 0
		x := pressedKeys[i:]

		for _, key := range []string{"222", "22", "2", "333", "33", "3", "444", "44", "4", "555", "55", "5", "666", "66", "6", "7777", "777", "77", "7", "888", "88", "8", "9999", "999", "99", "9"} {
			if strings.HasPrefix(x, key) {
				j := i + len(key)
				if j < len(dp) {
					res += dp[j]
				} else {
					res += 1
				}
				res %= mod
				if len(key) == 1 {
					break
				}
			}
		}
		dp[i] = res
	}

	return dp[0]
}
