package main

import (
	"sync"
	"time"
)

type slidingWindowLogEntry struct {
	Time int64
}

type SlidingWindowLog struct {
	mu       sync.Mutex
	Limit    int
	WindowMs int64
	Requests map[string][]slidingWindowLogEntry
}

func NewSlidingWindowLog(limit int, windowMs int64) *SlidingWindowLog {
	return &SlidingWindowLog{
		Limit:    limit,
		WindowMs: windowMs,
		Requests: make(map[string][]slidingWindowLogEntry),
	}
}

func (s *SlidingWindowLog) IsAllowed(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentTime := time.Now().UnixMilli()
	if _, ok := s.Requests[ip]; !ok {
		s.Requests[ip] = []slidingWindowLogEntry{{Time: currentTime}}
		return true
	}
	var valid []slidingWindowLogEntry
	for _, req := range s.Requests[ip] {
		if currentTime-req.Time < s.WindowMs {
			valid = append(valid, req)
		}
	}
	if len(valid) < s.Limit {
		s.Requests[ip] = append(valid, slidingWindowLogEntry{Time: currentTime})
		return true
	}
	return false
}
