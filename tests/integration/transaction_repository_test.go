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

func TestCreateTransaction(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	paymentMethod := entity.PaymentMethods["debit_card"]
	assert.Equal(t, err, nil)
	id := uuid.New().String()
	clientId := uuid.New().String()
	transaction, err := entity.NewTransaction(id, clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	repo := repository.NewTransactionRepository(DbConn)
	ctx := context.Background()
	err = repo.Save(ctx, transaction)
	assert.Equal(t, err, nil)
	transactionSaved, err := repo.GetByID(ctx, transaction.GetID())
	assert.Equal(t, err, nil)
	assert.Equal(t, transactionSaved.GetClientID(), transaction.GetClientID())
	assert.Equal(t, transactionSaved.GetValue(), transaction.GetValue())
	assert.Equal(t, transactionSaved.GetDescription(), transaction.GetDescription())
	assert.Equal(t, transactionSaved.GetPaymentMethod(), transaction.GetPaymentMethod())
	assert.Equal(t, transactionSaved.GetCreatedAt(), transaction.GetCreatedAt())
	assert.Equal(t, transactionSaved.GetCard().GetLastDigits(), transaction.GetCard().GetLastDigits())
	assert.Equal(t, transactionSaved.GetCard().GetVerificationCode(), transaction.GetCard().GetVerificationCode())
	assert.Equal(t, transactionSaved.GetCard().GetValidDate(), transaction.GetCard().GetValidDate())
	assert.Equal(t, transactionSaved.GetCard().GetOwnerName(), transaction.GetCard().GetOwnerName())
}
