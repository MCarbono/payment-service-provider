package integration

import (
	"context"
	"payment-service-provider/application/usecase"
	"payment-service-provider/domain/entity"
	db "payment-service-provider/infra/db/sqlc"
	"payment-service-provider/infra/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldListTransactions(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	assert.Equal(t, err, nil)
	clientId := uuid.New().String()
	transactionRepo := repository.NewTransactionRepository(DbConn)
	transactionDebitCard, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, entity.PaymentMethods["debit_card"])
	err = transactionRepo.Save(context.Background(), transactionDebitCard)
	transactionCreditCard, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, entity.PaymentMethods["credit_card"])
	err = transactionRepo.Save(context.Background(), transactionCreditCard)
	assert.Equal(t, err, nil)
	expected := []usecase.ListTransactionsOutput{
		{
			ID:                   transactionDebitCard.GetID(),
			ClientID:             transactionDebitCard.GetClientID(),
			Description:          transactionDebitCard.GetDescription(),
			Value:                transactionDebitCard.GetValue(),
			CardOwnerName:        transactionDebitCard.GetCard().GetOwnerName(),
			CardVerificationCode: transactionDebitCard.GetCard().GetVerificationCode(),
			CardLastDigits:       transactionDebitCard.GetCard().GetLastDigits(),
			CardValidDate:        transactionDebitCard.GetCard().GetValidDate(),
			PaymentMethod:        transactionDebitCard.GetPaymentMethod(),
			CreatedAt:            transactionDebitCard.GetCreatedAt(),
		},
		{
			ID:                   transactionCreditCard.GetID(),
			ClientID:             transactionCreditCard.GetClientID(),
			Description:          transactionCreditCard.GetDescription(),
			Value:                transactionCreditCard.GetValue(),
			CardOwnerName:        transactionCreditCard.GetCard().GetOwnerName(),
			CardVerificationCode: transactionCreditCard.GetCard().GetVerificationCode(),
			CardLastDigits:       transactionCreditCard.GetCard().GetLastDigits(),
			CardValidDate:        transactionCreditCard.GetCard().GetValidDate(),
			PaymentMethod:        transactionCreditCard.GetPaymentMethod(),
			CreatedAt:            transactionCreditCard.GetCreatedAt(),
		},
	}
	querier := db.New(DbConn)
	uc := usecase.NewListTransactions(querier)
	output, err := uc.Execute(context.Background(), &usecase.ListTransationsInput{ClientID: transactionDebitCard.GetClientID()})
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, output)
}

func TestShouldNotListTransactions(t *testing.T) {
	expected := []usecase.ListTransactionsOutput{}
	querier := db.New(DbConn)
	uc := usecase.NewListTransactions(querier)
	output, err := uc.Execute(context.Background(), &usecase.ListTransationsInput{ClientID: uuid.NewString()})
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, output)
}
