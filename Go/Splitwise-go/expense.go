package main

import (
	"fmt"
	"strings"
)

type Expense interface {
	Validate() error
	Apply(sheet *BalanceSheet) error
	GetID() int
}

type baseExpense struct {
	ID       int
	PaidBy   int
	Amount   int64 // paise
	Splits   []*Split
}

func (e *baseExpense) GetID() int { return e.ID }

func (e *baseExpense) Apply(sheet *BalanceSheet) error {
	for _, split := range e.Splits {
		if err := sheet.AddDebt(split.UserID, e.PaidBy, split.Amount); err != nil {
			return err
		}
	}
	return nil
}

type ExactExpense struct{ baseExpense }

func (e *ExactExpense) Validate() error {
	var total int64
	for _, s := range e.Splits {
		total += s.Amount
	}
	if total != e.Amount {
		return fmt.Errorf("total expense does not match amount")
	}
	return nil
}

type EqualExpense struct{ baseExpense }

func (e *EqualExpense) Validate() error {
	n := int64(len(e.Splits))
	base := e.Amount / n
	for _, s := range e.Splits {
		s.Amount = base
	}
	e.Splits[n-1].Amount += e.Amount - base*n
	return nil
}

type PercentageExpense struct{ baseExpense }

func (e *PercentageExpense) Validate() error {
	var totalPct int
	for _, s := range e.Splits {
		totalPct += s.Percentage
	}
	if totalPct != 100 {
		return fmt.Errorf("total percentage must be 100")
	}

	var allocated int64
	smallest := e.Amount
	smallestIndex := 0
	for i, s := range e.Splits {
		s.Amount = (e.Amount * int64(s.Percentage)) / 100
		allocated += s.Amount
		if s.Amount < smallest {
			smallest = s.Amount
			smallestIndex = i
		}
	}
	e.Splits[smallestIndex].Amount += e.Amount - allocated
	return nil
}

func CreateExpense(id int, typ string, paidBy int, amountPaise int64, splits []*Split) (Expense, error) {
	if paidBy == 0 || amountPaise == 0 || len(splits) == 0 {
		return nil, fmt.Errorf("invalid expense")
	}
	base := baseExpense{ID: id, PaidBy: paidBy, Amount: amountPaise, Splits: splits}
	switch strings.ToLower(typ) {
	case "exact":
		return &ExactExpense{base}, nil
	case "equal":
		return &EqualExpense{base}, nil
	case "percentage":
		return &PercentageExpense{base}, nil
	default:
		return nil, fmt.Errorf("invalid expense type")
	}
}
