package main

import "fmt"

type Split struct {
	UserID     int
	Amount     int64 // paise (filled by Equal/Percentage validate, or Exact input)
	Percentage int
}

func NewExactSplit(userID int, amountRupees float64) (*Split, error) {
	paise, err := ToAmount(amountRupees)
	if err != nil {
		return nil, err
	}
	return &Split{UserID: userID, Amount: paise}, nil
}

func NewPercentageSplit(userID int, percentage int) (*Split, error) {
	if percentage < 0 || percentage > 100 {
		return nil, fmt.Errorf("percentage must be 0–100")
	}
	return &Split{UserID: userID, Percentage: percentage}, nil
}

func NewEqualSplit(userID int) *Split {
	return &Split{UserID: userID}
}
