package main

import "fmt"

type PaymentStrategy interface {
	Pay(amount float64)
}

type CreditCardPayment struct {
	CardNumber string
}

func (c *CreditCardPayment) Pay(amount float64) {
	fmt.Printf("Paid %.2f using Credit Card [%s]\n", amount, c.CardNumber)
}

type PayPalPayment struct {
	Email string
}

func (p *PayPalPayment) Pay(amount float64) {
	fmt.Printf("Paid %.2f using PayPal account [%s]\n", amount, p.Email)
}

type PaymentProcessor struct {
	strategy PaymentStrategy
}

func (pp *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
	pp.strategy = strategy
}

func (pp *PaymentProcessor) Pay(amount float64) {
	if pp.strategy == nil {
		fmt.Println("No payment strategy selected!")
		return
	}
	pp.strategy.Pay(amount)
}

// Usage
func main() {
	processor := &PaymentProcessor{}

	// Use Credit Card
	processor.SetStrategy(&CreditCardPayment{CardNumber: "1234-5678-9012-3456"})
	processor.Pay(100.50)

	// Switch to PayPal
	processor.SetStrategy(&PayPalPayment{Email: "user@example.com"})
	processor.Pay(250.75)
}
