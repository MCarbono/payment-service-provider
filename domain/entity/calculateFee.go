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

func FeeCalculatorFactory(method paymentMethod) (Fee, error) {
	if method == debitCard {
		return newFeeWithDebitCard(), nil
	}
	if method == creditCard {
		return newFeeWithCreditCard(), nil
	}
	return nil, fmt.Errorf("invalid paymentMethod: %s", method)
}
