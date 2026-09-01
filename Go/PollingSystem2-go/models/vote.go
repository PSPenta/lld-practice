package models

import "time"

type Vote struct {
	PollID      int
	Option      string
	UserID      int
	SubmittedAt int64
}

func NewVote(pollID int, option string, userID int) (*Vote, error) {
	if pollID == 0 || option == "" || userID == 0 {
		return nil, ErrInvalidVote
	}
	return &Vote{
		PollID:      pollID,
		Option:      option,
		UserID:      userID,
		SubmittedAt: time.Now().UnixMilli(),
	}, nil
}
