package main

import "fmt"

type PaymentGateway struct {
	gateways map[PaymentMethod]BankGateway
}

func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{
		gateways: map[PaymentMethod]BankGateway{
			MethodUPI:     UPIGateway{},
			MethodCard:    CardGateway{},
			MethodNetBank: NetBankingGateway{},
		},
	}
}

func (pg *PaymentGateway) Process(payment *Payment) error {
	gateway, ok := pg.gateways[payment.Method]
	if !ok {
		payment.Status = StatusFailed
		return fmt.Errorf("unsupported payment method: %s", payment.Method)
	}
	if err := gateway.ProcessPayment(payment); err != nil {
		payment.Status = StatusFailed
		return err
	}
	return nil
}

func (pg *PaymentGateway) Refund(payment *Payment) error {
	if payment.Status != StatusSuccess {
		return fmt.Errorf("cannot refund non-successful payment")
	}
	payment.Status = StatusPending
	return nil
}
