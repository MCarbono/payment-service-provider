package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"payment-service-provider/application/usecase"
	"payment-service-provider/domain/entity"
	db "payment-service-provider/infra/db/sqlc"
	"payment-service-provider/infra/http/controllers"
	"payment-service-provider/infra/http/router"
	"payment-service-provider/infra/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListTransactions(t *testing.T) {
	transactionRepo := repository.NewTransactionRepository(DbConn)
	querier := db.New(DbConn)
	uc := usecase.NewListTransactions(querier)
	controller := controllers.NewTransaction(&usecase.ProcessTransaction{}, uc, &usecase.ClientBalance{})
	router := router.New(controller)
	t.Run("Should list transactions", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		assert.Equal(t, err, nil)
		clientId := uuid.New().String()
		transactionDebitCard, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, entity.PaymentMethods["debit_card"])
		err = transactionRepo.Save(context.Background(), transactionDebitCard)
		transactionCreditCard, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, entity.PaymentMethods["credit_card"])
		err = transactionRepo.Save(context.Background(), transactionCreditCard)
		assert.Equal(t, err, nil)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", clientId), nil)
		assert.Equal(t, err, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var output []usecase.ListTransactionsOutput
		err = json.Unmarshal(rec.Body.Bytes(), &output)
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
		assert.Equal(t, expected, output)
	})

	t.Run("Should List empty transactions", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", uuid.New().String()), nil)
		assert.Equal(t, err, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var output []usecase.ListTransactionsOutput
		err = json.Unmarshal(rec.Body.Bytes(), &output)
		expected := []usecase.ListTransactionsOutput{}
		assert.Equal(t, expected, output)
	})
}
