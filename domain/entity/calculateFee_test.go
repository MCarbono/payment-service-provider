package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFeeWithDebitCard(t *testing.T) {
	feeCalculator, err := FeeCalculatorFactory(debitCard)
	assert.Equal(t, err, nil)
	fee := feeCalculator.calculate(100)
	assert.Equal(t, fee, float32(3.0))
}

func TestCalculateFeeWithCreditCard(t *testing.T) {
	feeCalculator, err := FeeCalculatorFactory(creditCard)
	assert.Equal(t, err, nil)
	fee := feeCalculator.calculate(100)
	assert.Equal(t, fee, float32(5))
}

func TestFailInvalidMethod(t *testing.T) {
	_, err := FeeCalculatorFactory(paid)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid paymentMethod: paid")
}
