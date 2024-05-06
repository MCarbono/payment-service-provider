package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateCard(t *testing.T) {
	validDate := time.Now().AddDate(5, 0, 0)
	c, err := NewCard("Teste da Silva", "123", "1111-1111-1111-1111", validDate)
	assert.Nil(t, err)
	assert.Equal(t, "Teste da Silva", c.ownerName)
	assert.Equal(t, 4, len(c.GetLastDigits()))
	assert.Equal(t, "1111", c.lastDigits)
	assert.Equal(t, "123", c.verificationCode)
	assert.Equal(t, validDate.UTC().Truncate(time.Millisecond), c.validDate)
}

func TestShouldNotCreateCardInvalidVerificationCode(t *testing.T) {
	validDate := time.Now().AddDate(5, 0, 0)
	c, err := NewCard("Teste da Silva", "1234", "1111-1111-1111-1111", validDate)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid verificationCode length. Expected 3, got 4")
}

func TestShouldNotCreateCardInvalidCardNumber(t *testing.T) {
	validDate := time.Now().AddDate(5, 0, 0)
	c, err := NewCard("Teste da Silva", "123", "1111-1111-1111-11111111111", validDate)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid cardNumber length. Expected 16, got 23")
}

func TestShouldNotCreateCardInvalidOwnerName(t *testing.T) {
	validDate := time.Now().AddDate(5, 0, 0)
	c, err := NewCard("", "123", "1111-1111-1111-1111", validDate)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid ownerName length. Expected at least 8, got 0")
}

func TestShouldNotCreateCardBecauseIsExpired(t *testing.T) {
	validDate := time.Now().AddDate(0, 0, -1)
	c, err := NewCard("", "123", "1111-1111-1111-1111", validDate)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "card is expired")
}
