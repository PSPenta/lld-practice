package main

import (
	"fmt"
	"math"
	"strings"
)

type Expense struct {
	PaidBy string
	Amount float64
	Splits []Split
}

func (e *Expense) Apply(sheet *BalanceSheet) {
	for _, split := range e.Splits {
		if split.User != e.PaidBy {
			sheet.AddBalance(split.User, split.Amount)
		} else {
			sheet.AddBalance(split.User, -split.Amount)
		}
	}
}

type ExactExpense struct {
	Expense
}

func (e *ExactExpense) Validate() error {
	total := 0.0
	for _, split := range e.Splits {
		total += split.Amount
	}
	if total != e.Amount {
		return fmt.Errorf("total expense does not match amount")
	}
	return nil
}

type EqualExpense struct {
	Expense
}

func (e *EqualExpense) Validate() error {
	amount := math.Round(e.Amount/float64(len(e.Splits))*100) / 100
	for i := range e.Splits {
		e.Splits[i].Amount = amount
	}
	return nil
}

type PercentageExpense struct {
	Expense
}

func (e *PercentageExpense) Validate() error {
	totalPercentage := 0.0
	for _, split := range e.Splits {
		totalPercentage += split.Percentage
	}
	if totalPercentage != 100 {
		return fmt.Errorf("total percentage must be 100")
	}
	for i := range e.Splits {
		e.Splits[i].Amount = math.Round(e.Amount*e.Splits[i].Percentage/100*100) / 100
	}
	return nil
}

type expenseValidator interface {
	Validate() error
	Apply(sheet *BalanceSheet)
}

func CreateExpense(expenseType, paidBy string, amount float64, splits []Split) (expenseValidator, error) {
	switch strings.ToLower(expenseType) {
	case "exact":
		return &ExactExpense{Expense{PaidBy: paidBy, Amount: amount, Splits: splits}}, nil
	case "equal":
		return &EqualExpense{Expense{PaidBy: paidBy, Amount: amount, Splits: splits}}, nil
	case "percentage":
		return &PercentageExpense{Expense{PaidBy: paidBy, Amount: amount, Splits: splits}}, nil
	default:
		return nil, fmt.Errorf("invalid expense type")
	}
}
