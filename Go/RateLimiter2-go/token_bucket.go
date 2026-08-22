package main

import (
	"sync"
	"time"
)

type tokenBucketState struct {
	Tokens         int
	LastRefilledAt int64
}

type TokenBucketStrategy struct {
	mu         sync.Mutex
	Capacity   int
	RefillRate int
	Bucket     map[string]tokenBucketState
}

func NewTokenBucketStrategy(capacity, refillRate int) *TokenBucketStrategy {
	return &TokenBucketStrategy{
		Capacity:   capacity,
		RefillRate: refillRate,
		Bucket:     make(map[string]tokenBucketState),
	}
}

func (t *TokenBucketStrategy) refill(ip string) {
	currentTime := time.Now().UnixMilli()
	ipRequests := t.Bucket[ip]
	timeElapsed := float64(currentTime-ipRequests.LastRefilledAt) / 1000.0
	refills := int(timeElapsed * float64(t.RefillRate))
	if refills > 0 {
		t.Bucket[ip] = tokenBucketState{
			Tokens:         min(t.Capacity, ipRequests.Tokens+refills),
			LastRefilledAt: currentTime,
		}
	}
}

func (t *TokenBucketStrategy) IsAllowed(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Capacity == 0 {
		return false
	}
	if _, ok := t.Bucket[ip]; !ok {
		t.Bucket[ip] = tokenBucketState{
			Tokens:         t.Capacity - 1,
			LastRefilledAt: time.Now().UnixMilli(),
		}
		return true
	}
	if t.RefillRate > 0 {
		t.refill(ip)
	}
	ipRequests := t.Bucket[ip]
	if ipRequests.Tokens > 0 {
		t.Bucket[ip] = tokenBucketState{
			Tokens:         ipRequests.Tokens - 1,
			LastRefilledAt: ipRequests.LastRefilledAt,
		}
		return true
	}
	return false
}
