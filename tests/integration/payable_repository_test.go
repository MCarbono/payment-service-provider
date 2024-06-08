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
	payable, err := entity.PayableFactory(transaction)
	assert.Equal(t, err, nil)
	repo := repository.NewPayableRepository(DbConn)
	err = repo.Save(context.Background(), payable.GetData())
	assert.Equal(t, err, nil)
	savedPayable, err := repo.GetByID(context.Background(), payable.GetData().GetID())
	assert.Equal(t, err, nil)
	assert.NotNil(t, savedPayable)
	assert.Equal(t, savedPayable.GetClientID(), payable.GetData().GetClientID())
	assert.Equal(t, savedPayable.GetAmount(), payable.GetData().GetAmount())
	assert.Equal(t, savedPayable.GetFeeAmount(), payable.GetData().GetFeeAmount())
	assert.Equal(t, savedPayable.GetStatus(), payable.GetData().GetStatus())
	assert.Equal(t, savedPayable.GetTransactionID(), payable.GetData().GetTransactionID())
	assert.Equal(t, savedPayable.GetCreatedAt(), payable.GetData().GetCreatedAt())
}
