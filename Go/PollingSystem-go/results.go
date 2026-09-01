package main

import (
	"sync"
	"time"
)

var (
	resultsMu  sync.RWMutex
	allResults []Result
	pollsMu    sync.RWMutex
	allPolls   []Poll
)

func SubmitResult(result *Result) {
	if result == nil {
		return
	}
	resultsMu.Lock()
	defer resultsMu.Unlock()
	allResults = append(allResults, *result)
}

func GetStatistics(poll *Poll) map[string]any {
	resultsMu.RLock()
	defer resultsMu.RUnlock()
	var filtered []Result
	for _, res := range allResults {
		if res.PollID == poll.ID {
			filtered = append(filtered, res)
		}
	}

	ratings := make(map[string]int)
	for _, res := range filtered {
		ratings[res.Option]++
	}

	stats := make(map[string]float64)
	for _, option := range poll.Options {
		if len(filtered) > 0 {
			stats[option] = float64(ratings[option]) / float64(len(filtered))
		} else {
			stats[option] = 0
		}
	}

	return map[string]any{
		"question":   poll.Question,
		"statistics": stats,
	}
}

func AddPoll(poll *Poll) {
	if poll == nil {
		return
	}
	pollsMu.Lock()
	defer pollsMu.Unlock()
	allPolls = append(allPolls, *poll)
}

func GetActivePolls() []Poll {
	pollsMu.RLock()
	defer pollsMu.RUnlock()
	now := time.Now().UnixMilli()
	var active []Poll
	for _, poll := range allPolls {
		if poll.ValidTill.UnixMilli() > now {
			active = append(active, poll)
		}
	}
	return active
}

func GetCompletedPolls(adminID int) []Poll {
	pollsMu.RLock()
	defer pollsMu.RUnlock()
	now := time.Now().UnixMilli()
	var completed []Poll
	for _, poll := range allPolls {
		if poll.CreateBy == adminID && poll.ValidTill.UnixMilli() < now {
			completed = append(completed, poll)
		}
	}
	return completed
}
