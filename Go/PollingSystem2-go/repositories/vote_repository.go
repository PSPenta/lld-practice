package repositories

import (
	"fmt"
	"sync"

	"lld-practice/pollingsystem2-go/models"
)

type VoteRepository struct {
	mu    sync.RWMutex
	votes []*models.Vote
}

func NewVoteRepository() *VoteRepository {
	return &VoteRepository{votes: []*models.Vote{}}
}

func (r *VoteRepository) Add(vote *models.Vote) error {
	if vote == nil {
		return fmt.Errorf("invalid vote")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range r.votes {
		if v.PollID == vote.PollID && v.UserID == vote.UserID {
			return fmt.Errorf("user has already voted on this poll!")
		}
	}
	r.votes = append(r.votes, vote)
	return nil
}

func (r *VoteRepository) GetStatistics(poll *models.Poll) map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.Vote
	for _, v := range r.votes {
		if v.PollID == poll.ID {
			filtered = append(filtered, v)
		}
	}

	ratings := make(map[string]int)
	for _, v := range filtered {
		ratings[v.Option]++
	}

	total := len(filtered)
	stats := make(map[string]float64)
	counts := make(map[string]int)
	for _, option := range poll.Options {
		count := ratings[option]
		counts[option] = count
		if total == 0 {
			stats[option] = 0
		} else {
			stats[option] = float64(count) / float64(total)
		}
	}

	return map[string]any{
		"question":   poll.Question,
		"totalVotes": total,
		"counts":     counts,
		"statistics": stats,
	}
}
