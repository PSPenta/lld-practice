package main

import (
	"sync"
	"time"
)

type requestEntry struct {
	Path      string
	Timestamp int64
}

type FixedWindowRatelimiter struct {
	mu         sync.Mutex
	WindowMs   int64
	RateLimit  int
	RequestLog map[string][]requestEntry
}

func NewFixedWindowRatelimiter(windowMs int64, rateLimit int) *FixedWindowRatelimiter {
	return &FixedWindowRatelimiter{
		WindowMs:   windowMs,
		RateLimit:  rateLimit,
		RequestLog: make(map[string][]requestEntry),
	}
}

func (f *FixedWindowRatelimiter) IsRequestAllowed(userID, path string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	currentTime := time.Now().UnixMilli()

	existing, ok := f.RequestLog[userID]
	if !ok {
		f.RequestLog[userID] = []requestEntry{{Path: path, Timestamp: currentTime}}
		return true
	}

	var valid []requestEntry
	for _, req := range existing {
		if currentTime-req.Timestamp < f.WindowMs {
			valid = append(valid, req)
		}
	}

	if len(valid) >= f.RateLimit {
		return false
	}

	f.RequestLog[userID] = append(valid, requestEntry{Path: path, Timestamp: currentTime})
	return true
}
