package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewTransactionWithDebitCard(t *testing.T) {
	id := uuid.New().String()
	clientID := uuid.New().String()
	cardValidDate := time.Now().AddDate(5, 0, 0)
	card, err := NewCard("Teste da Silva", "123", "1234-1234-1234-1234", cardValidDate)
	assert.Nil(t, err)
	createdAt := time.Now().UTC()
	transaction, err := NewTransaction(id, clientID, "product", 100, createdAt, card, debitCard)
	assert.Nil(t, err)
	assert.Equal(t, transaction.id, id)
	assert.Equal(t, transaction.clientID, clientID)
	assert.Equal(t, transaction.description, "product")
	assert.Equal(t, transaction.paymentMethod, debitCard)
	assert.Equal(t, transaction.GetValue(), float32(100))
	assert.Equal(t, transaction.createdAt, createdAt.UTC().Truncate(time.Millisecond))
	assert.Equal(t, transaction.GetCard().ownerName, card.ownerName)
	assert.Equal(t, transaction.GetCard().lastDigits, card.lastDigits)
	assert.Equal(t, transaction.GetCard().verificationCode, card.verificationCode)
	assert.Equal(t, transaction.GetCard().validDate, card.validDate)
}

func TestCreateNewTransactionWithCreditCard(t *testing.T) {
	id := uuid.New().String()
	clientID := uuid.New().String()
	cardValidDate := time.Now().AddDate(5, 0, 0)
	card, err := NewCard("Teste da Silva", "123", "1234-1234-1234-1234", cardValidDate)
	assert.Nil(t, err)
	createdAt := time.Now().UTC()
	transaction, err := NewTransaction(id, clientID, "product", 100, createdAt, card, creditCard)
	assert.Nil(t, err)
	assert.Equal(t, transaction.id, id)
	assert.Equal(t, transaction.clientID, clientID)
	assert.Equal(t, transaction.description, "product")
	assert.Equal(t, transaction.paymentMethod, creditCard)
	assert.Equal(t, transaction.GetValue(), float32(100))
	assert.Equal(t, transaction.createdAt, createdAt.UTC().Truncate(time.Millisecond))
	assert.Equal(t, transaction.GetCard().ownerName, card.ownerName)
	assert.Equal(t, transaction.GetCard().lastDigits, card.lastDigits)
	assert.Equal(t, transaction.GetCard().verificationCode, card.verificationCode)
	assert.Equal(t, transaction.GetCard().validDate, card.validDate)
}

func TestShouldNotCreateTransactionInvalidClientID(t *testing.T) {
	id := uuid.New().String()
	createdAt := time.Now().UTC()
	transaction, err := NewTransaction(id, "123", "product", 100, createdAt, &Card{}, debitCard)
	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid clientID with format 123")
}
