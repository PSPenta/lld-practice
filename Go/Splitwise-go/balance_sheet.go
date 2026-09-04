package main

import "fmt"

type BalanceSheet struct {
	balances []*Balance
}

func NewBalanceSheet() *BalanceSheet {
	return &BalanceSheet{balances: make([]*Balance, 0)}
}

func (b *BalanceSheet) AddDebt(debtorID, creditorID int, amount int64) error {
	if debtorID == creditorID {
		return nil
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be a positive integer (paise)")
	}

	var existing *Balance
	for _, bal := range b.balances {
		if bal.DebtorID == debtorID && bal.CreditorID == creditorID {
			existing = bal
			break
		}
	}

	var opposite *Balance
	oppIdx := -1
	for i, bal := range b.balances {
		if bal.DebtorID == creditorID && bal.CreditorID == debtorID {
			opposite = bal
			oppIdx = i
			break
		}
	}

	if opposite != nil && opposite.Amount > 0 {
		if opposite.Amount > amount {
			opposite.Amount -= amount
			return nil
		}
		amount -= opposite.Amount
		b.balances = append(b.balances[:oppIdx], b.balances[oppIdx+1:]...)
	}

	if existing != nil {
		existing.Amount += amount
		return nil
	}
	if amount <= 0 {
		return nil
	}
	b.balances = append(b.balances, &Balance{
		DebtorID:   debtorID,
		CreditorID: creditorID,
		Amount:     amount,
	})
	return nil
}

func (b *BalanceSheet) GetPairwiseBalances() []*Balance {
	return b.balances
}

func (b *BalanceSheet) GetBalance(debtorID, creditorID int) *Balance {
	for _, bal := range b.balances {
		if bal.DebtorID == debtorID && bal.CreditorID == creditorID {
			return bal
		}
	}
	return nil
}
