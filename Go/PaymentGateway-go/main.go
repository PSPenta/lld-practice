package main

import "fmt"

func main() {
	gateway := NewPaymentGateway()

	p1, _ := NewPayment("PAY001", 500, MethodUPI, "user1")
	p2, _ := NewPayment("PAY002", 2500, MethodCard, "user2")
	p3, _ := NewPayment("PAY003", 15000, MethodNetBank, "user3")
	payments := []*Payment{p1, p2, p3}

	for _, payment := range payments {
		err := gateway.Process(payment)
		if err != nil {
			fmt.Printf("Payment %s failed via gateway: %v\n", payment.ID, err)
		} else {
			fmt.Printf("Payment %s processed successfully: status=%s amount=%.2f method=%s\n",
				payment.ID, payment.Status, payment.Amount, payment.Method)
		}
	}
}
