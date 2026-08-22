package main

import "fmt"

type Admin struct {
	ID int
}

func NewAdmin(id int) *Admin {
	return &Admin{ID: id}
}

func (a *Admin) CreatePoll(question string, options []string) (*Poll, error) {
	return NewPoll(question, options, a.ID)
}

func (a *Admin) ShowStatistics(poll *Poll) {
	statistics := GetStatistics(poll)
	fmt.Printf("Poll statistics for %d: %+v\n", poll.ID, statistics)
}
