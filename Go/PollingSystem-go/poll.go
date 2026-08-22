package main

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
	CreateBy      int
}

func NewPoll(question string, options []string, createdBy int) (*Poll, error) {
	if question == "" || len(options) <= 1 || createdBy == 0 {
		return nil, fmt.Errorf("invalid poll parameters!")
	}
	return &Poll{
		Question:      question,
		Options:       options,
		ScheduledTime: time.Now(),
		ValidTill:     time.Time{},
		CreateBy:      createdBy,
	}, nil
}
