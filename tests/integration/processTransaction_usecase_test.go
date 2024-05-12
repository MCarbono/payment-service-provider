package integration

import (
	"context"
	"payment-service-provider/application/usecase"
	"payment-service-provider/infra/repository"
	"payment-service-provider/tests/fake"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	transactionRepo := repository.NewTransactionRepository(DbConn)
	payableRepo := repository.NewPayableRepository(DbConn)
	t.Run("Should not process the transactions because the payment method is invalid", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		uc := usecase.NewProcessTransaction(Uow)
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
	})

	t.Run("Should be able to process a transaction with payment method as debit card", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		uc := usecase.NewProcessTransaction(Uow)
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
	})

	t.Run("Should be able to process a transaction with payment method as credit card", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		transactionRepo := repository.NewTransactionRepository(DbConn)
		payableRepo := repository.NewPayableRepository(DbConn)
		uc := usecase.NewProcessTransaction(Uow)
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
	})

	t.Run("Should not be able to process a transaction because something went wrong trying to save the payable", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		fake.NewFakeRepositoriesSavePayableError()
		defer fake.DestroyFakeRepositoriesImpl()
		uc := usecase.NewProcessTransaction(Uow)
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
		assert.NotNil(t, output)
		assert.Equal(t, err.Error(), "internal server error")
		_, err = transactionRepo.GetByID(context.Background(), output.TransactionDTO.ID)
		assert.Equal(t, err, repository.ErrTransactionNotFound)
		_, err = payableRepo.GetByID(context.Background(), output.PayableDTO.ID)
		assert.Equal(t, err, repository.ErrPayableNotFound)
	})

	t.Run("Should not be able to process a transaction because something went wrong trying to save the transaction", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		fake.NewFakeRepositoriesSaveTransactionError()
		defer fake.DestroyFakeRepositoriesImpl()
		uc := usecase.NewProcessTransaction(Uow)
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
		assert.NotNil(t, output)
		assert.Equal(t, err.Error(), "internal server error")
		_, err = transactionRepo.GetByID(context.Background(), output.TransactionDTO.ID)
		assert.Equal(t, err, repository.ErrTransactionNotFound)
		_, err = payableRepo.GetByID(context.Background(), output.PayableDTO.ID)
		assert.Equal(t, err, repository.ErrPayableNotFound)
	})

}
