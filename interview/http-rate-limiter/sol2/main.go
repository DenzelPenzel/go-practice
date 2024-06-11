package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens         int
	maxTokens      int
	refillRate     int
	refillInterval time.Duration
	mu             sync.Mutex
	cond           *sync.Cond
}

func NewRateLimiter(maxTokens, refillRate int, refillInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		refillInterval: refillInterval,
	}
	rl.cond = sync.NewCond(&rl.mu)
	go rl.startRefilling()
	return rl
}

func (rl *RateLimiter) startRefilling() {
	ticker := time.NewTicker(rl.refillInterval)

	for {
		<-ticker.C
		rl.mu.Lock()
		if rl.tokens < rl.maxTokens {
			rl.tokens = min(rl.tokens+rl.refillRate, rl.maxTokens)
			rl.cond.Broadcast() // Notify all waiting goroutines
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for rl.tokens <= 0 {
		rl.cond.Wait()
	}
	rl.tokens--
}

func main() {
	limiter := NewRateLimiter(5, 1, time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				fmt.Printf("Goroutine %d: Allowed\n", id)
			} else {
				fmt.Printf("Goroutine %d: Waiting\n", id)
				limiter.Wait()
				fmt.Printf("Goroutine %d: Proceeding after wait\n", id)
			}
		}(i)
	}

	wg.Wait()
}
