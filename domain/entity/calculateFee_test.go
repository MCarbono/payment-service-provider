package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFeeWithDebitCard(t *testing.T) {
	fee := calculateFee(100, debitCard.fee)
	assert.Equal(t, fee, float32(3.0))
}

func TestCalculateFeeWithCreditCard(t *testing.T) {
	fee := calculateFee(100, creditCard.fee)
	assert.Equal(t, fee, float32(5))
}
