package main

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	mu          sync.Mutex
	Limit       int
	WindowMs    int64
	Requests    map[string]int
	WindowStart int64
}

func NewFixedWindowCounter(limit int, windowMs int64) *FixedWindowCounter {
	return &FixedWindowCounter{
		Limit:       limit,
		WindowMs:    windowMs,
		Requests:    make(map[string]int),
		WindowStart: time.Now().UnixMilli(),
	}
}

func (f *FixedWindowCounter) IsAllowed(ip string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	currentTime := time.Now().UnixMilli()
	if currentTime-f.WindowStart >= f.WindowMs {
		f.WindowStart = currentTime
		f.Requests = make(map[string]int)
	}
	if _, ok := f.Requests[ip]; !ok {
		f.Requests[ip] = 1
		return true
	}
	if f.Requests[ip] >= f.Limit {
		return false
	}
	f.Requests[ip]++
	return true
}
