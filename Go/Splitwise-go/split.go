package main

import "fmt"

type Split struct {
	User       string
	Amount     float64
	Percentage float64
}

func NewSplit(user string, amount, percentage float64) (*Split, error) {
	if amount > 0 && percentage > 0 {
		return nil, fmt.Errorf("split cannot have both amount and percentage")
	}
	return &Split{User: user, Amount: amount, Percentage: percentage}, nil
}
