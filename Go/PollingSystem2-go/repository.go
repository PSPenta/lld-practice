package main

import (
	"fmt"
	"sync"
	"time"
)

type UserRepository struct {
	mu     sync.RWMutex
	users  []*User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: []*User{}, nextID: 1}
}

func (r *UserRepository) Create(email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("invalid user email!")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			return nil, fmt.Errorf("user already exists!")
		}
	}

	user, err := NewUser(r.nextID, email)
	if err != nil {
		return nil, err
	}
	r.nextID++
	r.users = append(r.users, user)
	return user, nil
}

type PollRepository struct {
	mu     sync.RWMutex
	polls  []*Poll
	nextID int
}

func NewPollRepository() *PollRepository {
	return &PollRepository{polls: []*Poll{}, nextID: 1}
}

func (r *PollRepository) NextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	r.nextID++
	return id
}

func (r *PollRepository) Add(poll *Poll) error {
	if poll == nil {
		return fmt.Errorf("invalid poll")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.polls {
		if p.ID == poll.ID {
			return fmt.Errorf("poll already exists!")
		}
	}
	r.polls = append(r.polls, poll)
	return nil
}

func (r *PollRepository) GetByID(id int) *Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, poll := range r.polls {
		if poll.ID == id {
			return poll
		}
	}
	return nil
}

func (r *PollRepository) Update(poll *Poll) error {
	if poll == nil {
		return fmt.Errorf("invalid poll!")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.polls {
		if p.ID == poll.ID {
			r.polls[i] = poll
			return nil
		}
	}
	return fmt.Errorf("poll not found!")
}

func (r *PollRepository) GetActive(now time.Time) []*Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var active []*Poll
	for _, poll := range r.polls {
		if !poll.IsCompleted(now) {
			active = append(active, poll)
		}
	}
	return active
}

func (r *PollRepository) GetCompletedByCreator(creatorID int, now time.Time) []*Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var completed []*Poll
	for _, poll := range r.polls {
		if poll.CreatedBy == creatorID && poll.IsCompleted(now) {
			completed = append(completed, poll)
		}
	}
	return completed
}

type VoteRepository struct {
	mu    sync.RWMutex
	votes []*Vote
}

func NewVoteRepository() *VoteRepository {
	return &VoteRepository{votes: []*Vote{}}
}

func (r *VoteRepository) Add(vote *Vote) error {
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

func (r *VoteRepository) GetStatistics(poll *Poll) map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*Vote
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
