/*
Rate Limiter question.

Take a continuous stream of timestamps on stdin
Example Input: 1074 1074 1076 1076 1076 1076

The input represents a timestamp at which a UDP packet was received from a
particular customer.

For each input, respond with either a 'd' which instructs the router to DROP the
packet, or an 'a' which instructs the router to ACCEPT the packet.

Accept as many packets as possible without exceeding:

3 packets per 1 second

10 packets per 5 seconds

1076 1074 1076 1076 1076 1076 1090
a     a     d   d    a    a    a
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type RateLimit struct {
	perSec     int
	perFiveSec int
	queue      []int
}

func NewRateLimiter(perSec, perFiveSec int) *RateLimit {
	return &RateLimit{
		queue:      []int{},
		perSec:     perSec,
		perFiveSec: perFiveSec,
	}
}

func (rl *RateLimit) shouldAccept(timestamp int, packetCount int) bool {
	rl.queue = append(rl.queue, timestamp)

	for len(rl.queue) > 0 && timestamp-rl.queue[0] > 5 {
		rl.queue = rl.queue[1:]
	}

	if len(rl.queue) <= rl.perFiveSec && packetCount < rl.perSec {
		return true
	}

	return false
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	limiter := NewRateLimiter(3, 10)

	scanner.Scan()
	arr := strings.Fields(scanner.Text())
	packetCount := 0
	res := make([]string, 0)

	for _, val := range arr {
		timestamp, err := strconv.Atoi(val)

		if err != nil {
			fmt.Println("Invalid timestamp:", err)
			continue
		}

		if limiter.shouldAccept(timestamp, packetCount) {
			res = append(res, "a")
			packetCount++
		} else {
			packetCount = int(math.Max(float64(0), float64(packetCount-1)))
			res = append(res, "b")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	fmt.Println(res)
}
