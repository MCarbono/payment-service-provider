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

func TestShouldGetClientBalanceAndTheStatusOfThePayableArePaid(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	paymentMethod := entity.PaymentMethods["debit_card"]
	assert.Equal(t, err, nil)
	clientId := uuid.New().String()
	firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	firstPayable, err := entity.GetPayable(uuid.NewString(), firstTransaction)
	assert.Equal(t, err, nil)
	secondPayable, err := entity.GetPayable(uuid.NewString(), secondTransaction)
	assert.Equal(t, err, nil)
	payableRepo := repository.NewPayableRepository(DbConn)
	err = payableRepo.Save(context.Background(), firstPayable)
	assert.Equal(t, err, nil)
	err = payableRepo.Save(context.Background(), secondPayable)
	assert.Equal(t, err, nil)
	querier := db.New(DbConn)
	uc := usecase.NewClientBalance(querier)
	output, err := uc.Execute(context.Background(), &usecase.ClientBalanceInput{ClientID: clientId})
	assert.Equal(t, err, nil)
	assert.Equal(t, output, &usecase.ClientBalanceOutput{
		Balance: usecase.ClientBalanceStatus{
			Paid:      194.0,
			WaitFunds: 0,
		},
	})
}

func TestShouldGetClientBalanceAndTheStatusOfThePayableAreWaitFunds(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	paymentMethod := entity.PaymentMethods["credit_card"]
	assert.Equal(t, err, nil)
	clientId := uuid.New().String()
	firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
	assert.Equal(t, err, nil)
	firstPayable, err := entity.GetPayable(uuid.NewString(), firstTransaction)
	assert.Equal(t, err, nil)
	secondPayable, err := entity.GetPayable(uuid.NewString(), secondTransaction)
	assert.Equal(t, err, nil)
	payableRepo := repository.NewPayableRepository(DbConn)
	err = payableRepo.Save(context.Background(), firstPayable)
	assert.Equal(t, err, nil)
	err = payableRepo.Save(context.Background(), secondPayable)
	assert.Equal(t, err, nil)
	querier := db.New(DbConn)
	uc := usecase.NewClientBalance(querier)
	output, err := uc.Execute(context.Background(), &usecase.ClientBalanceInput{ClientID: clientId})
	assert.Equal(t, err, nil)
	assert.Equal(t, output, &usecase.ClientBalanceOutput{
		Balance: usecase.ClientBalanceStatus{
			Paid:      0,
			WaitFunds: 190,
		},
	})
}

func TestShouldGetClientBalanceAndTheStatusOfThePayableAsPaidAndWaitFunds(t *testing.T) {
	defer cleanDatabaseTables(DbConn)
	card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
	assert.Equal(t, err, nil)
	clientId := uuid.New().String()
	firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 153.35, time.Now(), card, entity.PaymentMethods["credit_card"])
	assert.Equal(t, err, nil)
	secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 103.5, time.Now(), card, entity.PaymentMethods["debit_card"])
	assert.Equal(t, err, nil)
	firstPayable, err := entity.GetPayable(uuid.NewString(), firstTransaction)
	assert.Equal(t, err, nil)
	secondPayable, err := entity.GetPayable(uuid.NewString(), secondTransaction)
	assert.Equal(t, err, nil)
	payableRepo := repository.NewPayableRepository(DbConn)
	err = payableRepo.Save(context.Background(), firstPayable)
	assert.Equal(t, err, nil)
	err = payableRepo.Save(context.Background(), secondPayable)
	assert.Equal(t, err, nil)
	querier := db.New(DbConn)
	uc := usecase.NewClientBalance(querier)
	output, err := uc.Execute(context.Background(), &usecase.ClientBalanceInput{ClientID: clientId})
	assert.Equal(t, err, nil)
	assert.Equal(t, &usecase.ClientBalanceOutput{
		Balance: usecase.ClientBalanceStatus{
			Paid:      100,
			WaitFunds: 145,
		},
	}, output)
}
