package entity

import "fmt"

type Fee interface {
	calculate(value float32) float32
}

type FeeWithCreditCard struct {
	percentage float32
}

func newFeeWithCreditCard() *FeeWithCreditCard {
	return &FeeWithCreditCard{
		percentage: 0.05,
	}
}

func (f *FeeWithCreditCard) calculate(value float32) float32 {
	return value * f.percentage
}

type FeeWithDebitCard struct {
	percentage float32
}

func newFeeWithDebitCard() *FeeWithDebitCard {
	return &FeeWithDebitCard{
		percentage: 0.03,
	}
}

func (f *FeeWithDebitCard) calculate(value float32) float32 {
	return value * f.percentage
}

func FeeCalculatorFactory(method string) (Fee, error) {
	if method == "debit_card" {
		return newFeeWithDebitCard(), nil
	}
	if method == "credit_card" {
		return newFeeWithCreditCard(), nil
	}
	return nil, fmt.Errorf("invalid paymentMethod: %s", method)
}
