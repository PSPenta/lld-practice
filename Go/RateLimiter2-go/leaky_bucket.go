package main

import (
	"sync"
	"time"
)

type leakyBucketState struct {
	TotalRequests int
	LastLeakedAt  int64
}

type LeakyBucketStrategy struct {
	mu       sync.Mutex
	Capacity int
	LeakRate int
	Bucket   map[string]leakyBucketState
}

func NewLeakyBucketStrategy(capacity, leakRate int) *LeakyBucketStrategy {
	return &LeakyBucketStrategy{
		Capacity: capacity,
		LeakRate: leakRate,
		Bucket:   make(map[string]leakyBucketState),
	}
}

func (l *LeakyBucketStrategy) leak(ip string) {
	currentTime := time.Now().UnixMilli()
	ipRequests := l.Bucket[ip]
	timeElapsed := float64(currentTime-ipRequests.LastLeakedAt) / 1000.0
	leaks := int(float64(l.LeakRate) * timeElapsed)
	if leaks > 0 {
		l.Bucket[ip] = leakyBucketState{
			TotalRequests: max(0, ipRequests.TotalRequests-leaks),
			LastLeakedAt:  currentTime,
		}
	}
}

func (l *LeakyBucketStrategy) IsAllowed(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Capacity == 0 {
		return false
	}
	if _, ok := l.Bucket[ip]; !ok {
		l.Bucket[ip] = leakyBucketState{
			TotalRequests: 1,
			LastLeakedAt:  time.Now().UnixMilli(),
		}
		return true
	}
	if l.LeakRate > 0 {
		l.leak(ip)
	}
	ipRequests := l.Bucket[ip]
	if ipRequests.TotalRequests < l.Capacity {
		l.Bucket[ip] = leakyBucketState{
			TotalRequests: ipRequests.TotalRequests + 1,
			LastLeakedAt:  ipRequests.LastLeakedAt,
		}
		return true
	}
	return false
}
