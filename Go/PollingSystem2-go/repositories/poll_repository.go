package repositories

import (
	"fmt"
	"sync"
	"time"

	"lld-practice/pollingsystem2-go/models"
)

type PollRepository struct {
	mu     sync.RWMutex
	polls  []*models.Poll
	nextID int
}

func NewPollRepository() *PollRepository {
	return &PollRepository{polls: []*models.Poll{}, nextID: 1}
}

func (r *PollRepository) NextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	r.nextID++
	return id
}

func (r *PollRepository) Add(poll *models.Poll) error {
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

func (r *PollRepository) GetByID(id int) *models.Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, poll := range r.polls {
		if poll.ID == id {
			return poll
		}
	}
	return nil
}

func (r *PollRepository) Update(poll *models.Poll) error {
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

func (r *PollRepository) GetActive(now time.Time) []*models.Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var active []*models.Poll
	for _, poll := range r.polls {
		if !poll.IsCompleted(now) {
			active = append(active, poll)
		}
	}
	return active
}

func (r *PollRepository) GetCompletedByCreator(creatorID int, now time.Time) []*models.Poll {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var completed []*models.Poll
	for _, poll := range r.polls {
		if poll.CreatedBy == creatorID && poll.IsCompleted(now) {
			completed = append(completed, poll)
		}
	}
	return completed
}
