package main

import (
	"sync"
	"time"
)

const (
	tokenRateLimit           = 5
	tokenRefillRatePerSecond = 1
)

type TokenBucket struct {
	Limit      int
	Tokens     int
	RefillRate int
	LastRefill int64
}

func NewTokenBucket(limit, refillRate int) *TokenBucket {
	return &TokenBucket{
		Limit:      limit,
		Tokens:     limit,
		RefillRate: refillRate,
		LastRefill: time.Now().UnixMilli(),
	}
}

func (b *TokenBucket) refill() {
	now := time.Now().UnixMilli()
	timeElapsed := float64(now-b.LastRefill) / 1000.0
	tokensToAdd := int(timeElapsed * float64(b.RefillRate))
	b.Tokens = min(b.Limit, b.Tokens+tokensToAdd)
	b.LastRefill = now
}

func (b *TokenBucket) AllowRequest() bool {
	b.refill()
	if b.Tokens > 0 {
		b.Tokens--
		return true
	}
	return false
}

type TokenBucketLimiter struct {
	mu      sync.Mutex
	buckets map[string]*TokenBucket
}

func NewTokenBucketLimiter() *TokenBucketLimiter {
	return &TokenBucketLimiter{buckets: make(map[string]*TokenBucket)}
}

func (l *TokenBucketLimiter) AllowRequest(userIP string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.buckets[userIP]; !ok {
		l.buckets[userIP] = NewTokenBucket(tokenRateLimit, tokenRefillRatePerSecond)
	}
	return l.buckets[userIP].AllowRequest()
}
