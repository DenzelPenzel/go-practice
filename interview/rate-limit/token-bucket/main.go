package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate     float64       // Tokens added per second
	burst    float64       // Maximum number of tokens
	tokens   float64       // Current number of tokens
	last     time.Time     // Last time tokens were updated
	mutex    sync.Mutex    // Mutex to protect shared state
	cond     *sync.Cond    // Condition variable for waiting
	stopChan chan struct{} // Channel to signal stop
	stopped  bool          // Indicates if the limiter has been stopped
}

func NewRateLimiter(rate, burst float64) *RateLimiter {
	rl := &RateLimiter{
		rate:     rate,
		burst:    burst,
		tokens:   burst,
		last:     time.Now(),
		stopChan: make(chan struct{}),
	}

	rl.cond = sync.NewCond(&rl.mutex)
	go rl.refill()
	return rl
}

// refill continuously adds tokens to the bucket at the specified rate.
func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			if rl.stopped {
				rl.mutex.Unlock()
				return
			}
			now := time.Now()
			elapsed := now.Sub(rl.last).Seconds()
			rl.last = now
			rl.tokens += elapsed * rl.rate
			if rl.tokens > rl.burst {
				rl.tokens = rl.burst
			}
			rl.cond.Broadcast()
			rl.mutex.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

// Allow checks if a request can be processed immediately.
// It returns true if allowed, false otherwise
func () () {
	
}
 