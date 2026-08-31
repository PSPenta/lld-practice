package main

import "errors"

var (
	errInvalidUser = errors.New("invalid user parameters!")
	errInvalidVote = errors.New("invalid vote parameters!")
)
