package main

import (
	"fmt"
	"time"
)

type User struct {
	ID int
}

func NewUser(id int) *User {
	return &User{ID: id}
}

func (u *User) SubmitPoll(poll *Poll, option string) error {
	if poll == nil || option == "" {
		return fmt.Errorf("invalid poll or option")
	}
	if poll.ValidTill.Before(time.Now()) && !poll.ValidTill.IsZero() {
		return fmt.Errorf("poll has been expired!")
	}
	valid := false
	for _, o := range poll.Options {
		if o == option {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid option!")
	}
	SubmitResult(NewResult(poll.ID, option, u.ID))
	return nil
}
