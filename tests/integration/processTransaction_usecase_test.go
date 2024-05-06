package integration

import (
	"context"
	"payment-service-provider/application/usecase"
	"payment-service-provider/infra/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransactionInvalidPaymentMethod(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	transactionRepo := repository.NewTransactionRepository(DbConn)
	payableRepo := repository.NewPayableRepository(DbConn)
	uc := usecase.NewProcessTransaction(transactionRepo, payableRepo)
	input := &usecase.ProcessTransactionInput{
		ClientID:      uuid.New().String(),
		Value:         100.0,
		Description:   "description",
		PaymentMethod: "invalid_payment_method",
		Card: usecase.ProcessTransactionCardInput{
			Number:           "1111-1111-1111-1111",
			VerificationCode: "123",
			OwnerName:        "Teste da Silva",
			ValidDate:        time.Now().AddDate(5, 0, 0),
		},
	}
	output, err := uc.Execute(context.Background(), input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid payment method invalid_payment_method")
}

func TestShouldProcessTransactionSuccessfully(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	transactionRepo := repository.NewTransactionRepository(DbConn)
	payableRepo := repository.NewPayableRepository(DbConn)
	uc := usecase.NewProcessTransaction(transactionRepo, payableRepo)
	input := &usecase.ProcessTransactionInput{
		ClientID:      uuid.New().String(),
		Value:         100.0,
		Description:   "description",
		PaymentMethod: "debit_card",
		Card: usecase.ProcessTransactionCardInput{
			Number:           "1111-1111-1111-1111",
			VerificationCode: "123",
			OwnerName:        "Teste da Silva",
			ValidDate:        time.Now().AddDate(5, 0, 0),
		},
	}
	output, err := uc.Execute(context.Background(), input)
	assert.Nil(t, err)
	savedTransaction, err := transactionRepo.GetByID(context.Background(), output.TransactionDTO.ID)
	assert.Nil(t, err)
	assert.NotNil(t, savedTransaction)
	savedPayable, err := payableRepo.GetByID(context.Background(), output.PayableDTO.ID)
	assert.Nil(t, err)
	assert.NotNil(t, savedPayable)
	savedOutput := usecase.NewProcessTransactionDTO(savedTransaction, savedPayable)
	assert.Equal(t, savedOutput, output)
}

func TestShouldProcessTransactionSuccessfullyCreditCard(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	transactionRepo := repository.NewTransactionRepository(DbConn)
	payableRepo := repository.NewPayableRepository(DbConn)
	uc := usecase.NewProcessTransaction(transactionRepo, payableRepo)
	input := &usecase.ProcessTransactionInput{
		ClientID:      uuid.New().String(),
		Value:         100.0,
		Description:   "description",
		PaymentMethod: "credit_card",
		Card: usecase.ProcessTransactionCardInput{
			Number:           "1111-1111-1111-1111",
			VerificationCode: "123",
			OwnerName:        "Teste da Silva",
			ValidDate:        time.Now().AddDate(5, 0, 0),
		},
	}
	output, err := uc.Execute(context.Background(), input)
	assert.Nil(t, err)
	savedTransaction, err := transactionRepo.GetByID(context.Background(), output.TransactionDTO.ID)
	assert.Nil(t, err)
	assert.NotNil(t, savedTransaction)
	savedPayable, err := payableRepo.GetByID(context.Background(), output.PayableDTO.ID)
	assert.Nil(t, err)
	assert.NotNil(t, savedPayable)
	savedOutput := usecase.NewProcessTransactionDTO(savedTransaction, savedPayable)
	assert.Equal(t, savedOutput, output)
}
