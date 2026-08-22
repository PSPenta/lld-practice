package main

type RateLimiterStrategy interface {
	IsAllowed(ip string) bool
}

type RateLimiter struct {
	strategy RateLimiterStrategy
}

func NewRateLimiter(strategy RateLimiterStrategy) *RateLimiter {
	return &RateLimiter{strategy: strategy}
}

func (r *RateLimiter) IsAllowed(ip string) bool {
	if r.strategy == nil {
		panic("please set a strategy!")
	}
	return r.strategy.IsAllowed(ip)
}
