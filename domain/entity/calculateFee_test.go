package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFeeWithDebitCard(t *testing.T) {
	feeCalculator, err := FeeCalculatorFactory("debit_card")
	assert.Equal(t, err, nil)
	fee := feeCalculator.calculate(100)
	assert.Equal(t, fee, float32(3.0))
}

func TestCalculateFeeWithCreditCard(t *testing.T) {
	feeCalculator, err := FeeCalculatorFactory("credit_card")
	assert.Equal(t, err, nil)
	fee := feeCalculator.calculate(100)
	assert.Equal(t, fee, float32(5))
}

func TestFailInvalidMethod(t *testing.T) {
	_, err := FeeCalculatorFactory("invalid_method")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid paymentMethod: invalid_method")
}
