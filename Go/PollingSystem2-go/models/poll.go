package models

import (
	"fmt"
	"time"
)

type Poll struct {
	ID            int
	Question      string
	Options       []string
	ScheduledTime time.Time
	ValidTill     time.Time
	CreatedBy     int
	AssignedUsers []int
	IsPrivate     bool
	IsClosed      bool
}

func NewPoll(id int, question string, options []string, createdBy int, isPrivate, isClosed bool, duration time.Duration) (*Poll, error) {
	if id == 0 || question == "" || len(options) <= 1 || createdBy == 0 {
		return nil, fmt.Errorf("invalid poll parameters!")
	}

	seen := make(map[string]struct{}, len(options))
	for _, opt := range options {
		if _, ok := seen[opt]; ok {
			return nil, fmt.Errorf("poll options must be unique!")
		}
		seen[opt] = struct{}{}
	}

	if duration <= 0 {
		duration = 24 * time.Hour
	}

	now := time.Now()
	return &Poll{
		ID:            id,
		Question:      question,
		Options:       options,
		ScheduledTime: now,
		ValidTill:     now.Add(duration),
		CreatedBy:     createdBy,
		AssignedUsers: []int{},
		IsPrivate:     isPrivate,
		IsClosed:      isClosed,
	}, nil
}

func (p *Poll) IsCreator(userID int) bool {
	return p.CreatedBy == userID
}

func (p *Poll) IsAssigned(userID int) bool {
	for _, id := range p.AssignedUsers {
		if id == userID {
			return true
		}
	}
	return false
}

func (p *Poll) AssignVoter(userID int) error {
	if userID == 0 {
		return fmt.Errorf("invalid user!")
	}
	if userID == p.CreatedBy {
		return fmt.Errorf("you cannot assign yourself to your own poll!")
	}
	if p.IsAssigned(userID) {
		return fmt.Errorf("user already assigned to this poll!")
	}
	p.AssignedUsers = append(p.AssignedUsers, userID)
	return nil
}

func (p *Poll) IsCompleted(now time.Time) bool {
	if p.ValidTill.Before(now) {
		p.IsClosed = true
	}
	return p.IsClosed
}
