package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	id            string
	clientID      string
	description   string
	value         float32
	card          *Card
	paymentMethod paymentMethod
	createdAt     time.Time
}

func NewTransaction(
	id, clientID string, description string, value float32,
	createdAt time.Time, card *Card, paymentMethod paymentMethod,
) (*Transaction, error) {
	if err := uuid.Validate(clientID); err != nil {
		return nil, fmt.Errorf("invalid clientID with format %s", clientID)
	}
	t := &Transaction{
		id:            id,
		clientID:      clientID,
		description:   description,
		value:         value,
		card:          card,
		paymentMethod: paymentMethod,
		createdAt:     createdAt.UTC().Truncate(time.Millisecond),
	}
	return t, nil
}

func (t *Transaction) GetID() string {
	return t.id
}

func (t *Transaction) GetClientID() string {
	return t.clientID
}

func (t *Transaction) GetDescription() string {
	return t.description
}

func (t *Transaction) GetValue() float32 {
	return t.value
}

func (t *Transaction) GetCard() *Card {
	return t.card
}

func (t *Transaction) PaymentMethod() paymentMethod {
	return t.paymentMethod
}

func (t *Transaction) GetCreatedAt() time.Time {
	return t.createdAt
}
