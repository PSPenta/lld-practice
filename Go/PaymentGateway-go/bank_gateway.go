package main

import "fmt"

type BankGateway interface {
	ProcessPayment(payment *Payment) error
	Name() string
}

type UPIGateway struct{}

func (g UPIGateway) Name() string { return "UPI Gateway" }

func (g UPIGateway) ProcessPayment(payment *Payment) error {
	if payment.Amount > 100000 {
		return fmt.Errorf("UPI limit exceeded")
	}
	payment.Status = StatusSuccess
	return nil
}

type CardGateway struct{}

func (g CardGateway) Name() string { return "Card Gateway" }

func (g CardGateway) ProcessPayment(payment *Payment) error {
	if payment.Amount <= 0 {
		return fmt.Errorf("invalid card payment amount")
	}
	payment.Status = StatusSuccess
	return nil
}

type NetBankingGateway struct{}

func (g NetBankingGateway) Name() string { return "Net Banking Gateway" }

func (g NetBankingGateway) ProcessPayment(payment *Payment) error {
	payment.Status = StatusSuccess
	return nil
}

func GatewayForMethod(method PaymentMethod) BankGateway {
	switch method {
	case MethodUPI:
		return UPIGateway{}
	case MethodCard:
		return CardGateway{}
	case MethodNetBank:
		return NetBankingGateway{}
	default:
		return nil
	}
}
