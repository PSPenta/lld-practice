package main

import (
	"sync"
	"time"
)

type slidingWindowState struct {
	PreviousStart int64
	PreviousCount int
	CurrentStart  int64
	CurrentCount  int
}

type SlidingWindowCounter struct {
	mu       sync.Mutex
	Limit    int
	WindowMs int64
	Requests map[string]slidingWindowState
}

func NewSlidingWindowCounter(limit int, windowMs int64) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		Limit:    limit,
		WindowMs: windowMs,
		Requests: make(map[string]slidingWindowState),
	}
}

func (s *SlidingWindowCounter) IsAllowed(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentTime := time.Now().UnixMilli()

	window, ok := s.Requests[ip]
	if !ok {
		s.Requests[ip] = slidingWindowState{
			PreviousStart: currentTime - s.WindowMs,
			CurrentStart:  currentTime,
		}
		window = s.Requests[ip]
	}

	timeLapsed := currentTime - window.CurrentStart
	if timeLapsed >= s.WindowMs {
		window.PreviousStart = window.CurrentStart
		if timeLapsed >= s.WindowMs*2 {
			window.PreviousCount = 0
		} else {
			window.PreviousCount = window.CurrentCount
		}
		window.CurrentStart = currentTime
		window.CurrentCount = 0
		s.Requests[ip] = window
		timeLapsed = 0
	}

	ratio := float64(timeLapsed) / float64(s.WindowMs)
	effectivePreviousCount := float64(window.PreviousCount) * (1 - ratio)
	effectiveTotal := float64(window.CurrentCount) + effectivePreviousCount
	if effectiveTotal < float64(s.Limit) {
		window.CurrentCount++
		s.Requests[ip] = window
		return true
	}
	return false
}
