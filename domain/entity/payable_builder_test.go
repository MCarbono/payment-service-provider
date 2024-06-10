package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPayableBuilder(t *testing.T) {
	t.Run("Should create a payable object with method payment as debit card", func(t *testing.T) {
		clientId := uuid.New().String()
		card, err := NewCard("test owner name", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		assert.Equal(t, err, nil)
		transaction, err := NewTransaction(uuid.NewString(), clientId, "transaction Test", 100, time.Now(), card, debitCard)
		assert.Equal(t, err, nil)
		builder, err := BuilderFactory(transaction)
		assert.Equal(t, err, nil)
		director := GetPayableDirector()
		director.SetBuilder(builder)
		director.Construct(transaction)
		payable := builder.Build()
		assert.Equal(t, err, nil)
		assert.NotEmpty(t, payable.id)
		parsedUUID, err := uuid.Parse(payable.id)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(parsedUUID), 16)
		assert.Equal(t, payable.clientID, transaction.clientID)
		assert.Equal(t, payable.status, paid)
		assert.Equal(t, payable.feeAmount, float32(3))
		assert.Equal(t, payable.amount, float32(97))
		assert.Equal(t, payable.createdAt, transaction.createdAt)
		assert.Equal(t, payable.paymentDate, transaction.createdAt)
	})
	t.Run("Should create a payable object with method payment as credit card", func(t *testing.T) {
		clientId := uuid.New().String()
		card, err := NewCard("test owner name", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		assert.Equal(t, err, nil)
		transaction, err := NewTransaction(uuid.NewString(), clientId, "transaction Test", 100, time.Now(), card, creditCard)
		assert.Equal(t, err, nil)
		builder, err := BuilderFactory(transaction)
		assert.Equal(t, err, nil)
		director := GetPayableDirector()
		director.SetBuilder(builder)
		director.Construct(transaction)
		payable := builder.Build()
		assert.Equal(t, err, nil)
		assert.NotEmpty(t, payable.id)
		parsedUUID, err := uuid.Parse(payable.id)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(parsedUUID), 16)
		assert.Equal(t, payable.clientID, transaction.clientID)
		assert.Equal(t, payable.status, waitFunds)
		assert.Equal(t, payable.feeAmount, float32(5))
		assert.Equal(t, payable.amount, float32(95))
		assert.Equal(t, payable.createdAt, transaction.createdAt)
		assert.Equal(t, payable.paymentDate, transaction.createdAt.AddDate(0, 0, 30))
	})
}
