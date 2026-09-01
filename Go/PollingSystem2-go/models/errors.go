package models

import "errors"

var (
	ErrInvalidUser = errors.New("invalid user parameters!")
	ErrInvalidVote = errors.New("invalid vote parameters!")
)
