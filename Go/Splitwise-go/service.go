package main

import (
	"fmt"
	"strings"
)

type SplitwiseService struct {
	userID       int
	expenseID    int
	users        map[int]*User
	expenses     map[int]Expense
	balanceSheet *BalanceSheet
}

func NewSplitwiseService() *SplitwiseService {
	return &SplitwiseService{
		users:        make(map[int]*User),
		expenses:     make(map[int]Expense),
		balanceSheet: NewBalanceSheet(),
	}
}

func (s *SplitwiseService) AddUser(name, email string) (*User, error) {
	if name == "" || email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, fmt.Errorf("invalid user input")
	}
	for _, u := range s.users {
		if u.Email == email {
			return nil, fmt.Errorf("user already exists")
		}
	}
	s.userID++
	user, err := NewUser(s.userID, name, email)
	if err != nil {
		s.userID--
		return nil, err
	}
	s.users[s.userID] = user
	return user, nil
}

func (s *SplitwiseService) GetUser(id int) (*User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

type AddExpenseInput struct {
	Type    string
	PaidBy  int
	Amount  float64 // rupees at API boundary
	Splits  []*Split
}

func (s *SplitwiseService) AddExpense(input AddExpenseInput) (Expense, error) {
	if input.Type == "" || input.PaidBy == 0 || input.Amount == 0 || len(input.Splits) == 0 {
		return nil, fmt.Errorf("invalid expense input")
	}
	if _, ok := s.users[input.PaidBy]; !ok {
		return nil, fmt.Errorf("user not found")
	}
	for _, split := range input.Splits {
		if _, ok := s.users[split.UserID]; !ok {
			return nil, fmt.Errorf("user not found")
		}
	}

	amountPaise, err := ToAmount(input.Amount)
	if err != nil {
		return nil, err
	}

	s.expenseID++
	expense, err := CreateExpense(s.expenseID, input.Type, input.PaidBy, amountPaise, input.Splits)
	if err != nil {
		s.expenseID--
		return nil, err
	}
	if err := expense.Validate(); err != nil {
		s.expenseID--
		return nil, err
	}
	if err := expense.Apply(s.balanceSheet); err != nil {
		s.expenseID--
		return nil, err
	}
	s.expenses[s.expenseID] = expense
	return expense, nil
}

func (s *SplitwiseService) GetPairwiseBalances() []string {
	rows := s.balanceSheet.GetPairwiseBalances()
	if len(rows) == 0 {
		return []string{"All balances are settled!"}
	}
	out := make([]string, 0, len(rows))
	for _, bal := range rows {
		debtor := s.users[bal.DebtorID]
		creditor := s.users[bal.CreditorID]
		out = append(out, fmt.Sprintf(
			"User %s (%d) owes User %s (%d) %v",
			debtor.Email, bal.DebtorID, creditor.Email, bal.CreditorID, FromAmount(bal.Amount),
		))
	}
	return out
}

// SettleUp clears or partially clears debt. amountRupees nil => full settle.
func (s *SplitwiseService) SettleUp(payerID, payeeID int, amountRupees *float64) error {
	payer, ok1 := s.users[payerID]
	payee, ok2 := s.users[payeeID]
	if !ok1 || !ok2 {
		return fmt.Errorf("user not found")
	}
	bal := s.balanceSheet.GetBalance(payerID, payeeID)
	if bal == nil || bal.Amount <= 0 {
		return fmt.Errorf("balance not found")
	}

	var pay int64
	if amountRupees == nil {
		pay = bal.Amount
	} else {
		var err error
		pay, err = ToAmount(*amountRupees)
		if err != nil {
			return err
		}
	}
	if pay <= 0 {
		return fmt.Errorf("invalid settle amount")
	}
	if pay > bal.Amount {
		return fmt.Errorf("amount is greater than balance")
	}
	if err := s.balanceSheet.AddDebt(payeeID, payerID, pay); err != nil {
		return err
	}
	fmt.Printf("User %s (%d) settled up with User %s (%d) %v\n",
		payer.Email, payerID, payee.Email, payeeID, FromAmount(pay))
	return nil
}
