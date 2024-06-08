package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateAPayableWithDebitCard(t *testing.T) {
	clientId := uuid.New().String()
	card, err := NewCard("test owner name", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	assert.Equal(t, err, nil)
	transaction, err := NewTransaction(uuid.NewString(), clientId, "transaction Test", 100, time.Now(), card, debitCard)
	assert.Equal(t, err, nil)
	payable, err := PayableFactory(transaction)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, payable.GetData().id)
	parsedUUID, err := uuid.Parse(payable.GetData().id)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(parsedUUID), 16)
	assert.Equal(t, payable.GetData().clientID, transaction.clientID)
	assert.Equal(t, payable.GetData().status, paid)
	assert.Equal(t, payable.GetData().feeAmount, float32(3))
	assert.Equal(t, payable.GetData().amount, float32(97))
	assert.Equal(t, payable.GetData().createdAt, transaction.createdAt)
	assert.Equal(t, payable.GetData().paymentDate, transaction.createdAt)
}

func TestShouldCreateAPayableWithCreditCard(t *testing.T) {
	clientId := uuid.New().String()
	card, err := NewCard("test owner name", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	assert.Equal(t, err, nil)
	transaction, err := NewTransaction(uuid.NewString(), clientId, "transaction Test", 100, time.Now(), card, creditCard)
	assert.Equal(t, err, nil)
	payable, err := PayableFactory(transaction)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, payable.GetData().id)
	parsedUUID, err := uuid.Parse(payable.GetData().id)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(parsedUUID), 16)
	assert.Equal(t, payable.GetData().clientID, transaction.clientID)
	assert.Equal(t, payable.GetData().status, waitFunds)
	assert.Equal(t, payable.GetData().feeAmount, float32(5))
	assert.Equal(t, payable.GetData().amount, float32(95))
	assert.Equal(t, payable.GetData().createdAt, transaction.createdAt)
	assert.Equal(t, payable.GetData().paymentDate, transaction.createdAt.AddDate(0, 0, 30))
}
