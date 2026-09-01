package main

import "time"

type Result struct {
	PollID      int
	Option      string
	UserID      int
	SubmittedAt int64
}

func NewResult(pollID int, option string, userID int) *Result {
	return &Result{
		PollID:      pollID,
		Option:      option,
		UserID:      userID,
		SubmittedAt: time.Now().UnixMilli(),
	}
}
