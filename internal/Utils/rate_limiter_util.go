package utils

import (
	"sync"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*rate.Limiter // clientid vs limiter
	burst    int
	limit    rate.Limit
}

func NewRateLimiter(limitPerSecond int, burst int) *rateLimiter {

	return &rateLimiter{
		mu:       sync.Mutex{},
		visitors: make(map[string]*rate.Limiter),
		limit:    rate.Limit(limitPerSecond),
		burst:    burst,
	}
}

func (rl *rateLimiter) Allow(clientId string) bool {
	limiter := rl.GetLimiterForClient(clientId)
	return limiter.Allow()
}

func (rl *rateLimiter) GetLimiterForClient(clientId string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	limiter, exist := rl.visitors[clientId]
	if exist {
		return limiter
	} else {
		limiter = rate.NewLimiter(rl.limit, rl.burst)
		rl.visitors[clientId] = limiter
		return limiter
	}
}
