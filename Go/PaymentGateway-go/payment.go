package main

import "fmt"

type PaymentStatus string

const (
	StatusPending PaymentStatus = "PENDING"
	StatusSuccess PaymentStatus = "SUCCESS"
	StatusFailed  PaymentStatus = "FAILED"
)

type PaymentMethod string

const (
	MethodUPI     PaymentMethod = "UPI"
	MethodCard    PaymentMethod = "CARD"
	MethodNetBank PaymentMethod = "NET_BANKING"
)

type Payment struct {
	ID     string
	Amount float64
	Method PaymentMethod
	Status PaymentStatus
	UserID string
}

func NewPayment(id string, amount float64, method PaymentMethod, userID string) (*Payment, error) {
	if id == "" || amount <= 0 || userID == "" {
		return nil, fmt.Errorf("invalid payment details")
	}
	return &Payment{
		ID:     id,
		Amount: amount,
		Method: method,
		Status: StatusPending,
		UserID: userID,
	}, nil
}
