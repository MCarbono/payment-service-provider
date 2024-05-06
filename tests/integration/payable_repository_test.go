package integration

import (
	"context"
	"payment-service-provider/domain/entity"
	"payment-service-provider/infra/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayable(t *testing.T) {
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	paymentMethod := entity.PaymentMethods["debit_card"]
	assert.Equal(t, err, nil)
	id := uuid.New().String()
	clientId := uuid.New().String()
	transaction, err := entity.NewTransaction(id, clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	payable, err := entity.GetPayable(uuid.New().String(), transaction)
	assert.Equal(t, err, nil)
	repo := repository.NewPayableRepository(DbConn)
	err = repo.Save(context.Background(), payable)
	assert.Equal(t, err, nil)
	savedPayable, err := repo.GetByID(context.Background(), payable.GetID())
	assert.Equal(t, err, nil)
	assert.NotNil(t, savedPayable)
	assert.Equal(t, savedPayable.GetClientID(), payable.GetClientID())
	assert.Equal(t, savedPayable.GetAmount(), payable.GetAmount())
	assert.Equal(t, savedPayable.GetFeeAmount(), payable.GetFeeAmount())
	assert.Equal(t, savedPayable.GetStatus(), payable.GetStatus())
	assert.Equal(t, savedPayable.GetTransactionID(), payable.GetTransactionID())
	assert.Equal(t, savedPayable.GetCreatedAt(), payable.GetCreatedAt())

}
