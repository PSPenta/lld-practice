package main

import (
	"sync"
	"time"
)

const (
	leakyRateLimit    = 5
	leakyLeakRateInMs = 1000
)

type LeakyBucket struct {
	Limit    int
	LeakRate int64
	Requests int
	LastLeak int64
}

func NewLeakyBucket(limit int, leakRate int64) *LeakyBucket {
	return &LeakyBucket{Limit: limit, LeakRate: leakRate}
}

func (b *LeakyBucket) leak() {
	now := time.Now().UnixMilli()
	timeElapsed := now - b.LastLeak
	leakedRequests := int(timeElapsed / b.LeakRate)
	b.Requests = max(0, b.Requests-leakedRequests)
	b.LastLeak = now
}

func (b *LeakyBucket) AllowRequest() bool {
	b.leak()
	if b.Requests < b.Limit {
		b.Requests++
		return true
	}
	return false
}

type LeakyBucketLimiter struct {
	mu      sync.Mutex
	buckets map[string]*LeakyBucket
}

func NewLeakyBucketLimiter() *LeakyBucketLimiter {
	return &LeakyBucketLimiter{buckets: make(map[string]*LeakyBucket)}
}

func (l *LeakyBucketLimiter) AllowRequest(userIP string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.buckets[userIP]; !ok {
		l.buckets[userIP] = NewLeakyBucket(leakyRateLimit, leakyLeakRateInMs)
	}
	return l.buckets[userIP].AllowRequest()
}
