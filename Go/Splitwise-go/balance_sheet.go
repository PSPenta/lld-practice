package main

type BalanceSheet struct {
	balances map[string]float64
}

func NewBalanceSheet() *BalanceSheet {
	return &BalanceSheet{balances: make(map[string]float64)}
}

func (b *BalanceSheet) AddBalance(user string, amount float64) {
	b.balances[user] += amount
}

func (b *BalanceSheet) GetBalance(user string) float64 {
	return b.balances[user]
}

func (b *BalanceSheet) GetBalances() map[string]float64 {
	return b.balances
}
