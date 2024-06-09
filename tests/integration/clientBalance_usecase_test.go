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

func TestClientBalance(t *testing.T) {
	payableRepo := repository.NewPayableRepository(DbConn)
	querier := db.New(DbConn)
	uc := usecase.NewClientBalance(querier)
	c := controllers.NewTransaction(&usecase.ProcessTransaction{}, &usecase.ListTransactions{}, uc)
	r := router.New(c)
	t.Run("Should get client balance that returns only payables with status 'available'", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		paymentMethod := entity.PaymentMethods["debit_card"]
		assert.Equal(t, err, nil)
		clientId := uuid.New().String()
		firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
		assert.Equal(t, err, nil)
		secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
		assert.Equal(t, err, nil)
		firstPayable, err := entity.PayableFactory(firstTransaction)
		assert.Equal(t, err, nil)
		secondPayable, err := entity.PayableFactory(secondTransaction)
		assert.Equal(t, err, nil)
		err = payableRepo.Save(context.Background(), firstPayable)
		assert.Equal(t, err, nil)
		err = payableRepo.Save(context.Background(), secondPayable)
		assert.Equal(t, err, nil)
		url := fmt.Sprintf("/balance/%s", clientId)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Equal(t, err, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var output usecase.ClientBalanceOutput
		err = json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Equal(t, err, nil)
		assert.Equal(t, output, usecase.ClientBalanceOutput{
			Balance: usecase.ClientBalanceStatus{
				Paid:      194.0,
				WaitFunds: 0,
			},
		})
	})

	t.Run("Should get client balance that returns only payables with status 'wait_funds'", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		paymentMethod := entity.PaymentMethods["credit_card"]
		assert.Equal(t, err, nil)
		clientId := uuid.New().String()
		firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
		assert.Equal(t, err, nil)
		secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 100.0, time.Now(), card, paymentMethod)
		assert.Equal(t, err, nil)
		firstPayable, err := entity.PayableFactory(firstTransaction)
		assert.Equal(t, err, nil)
		secondPayable, err := entity.PayableFactory(secondTransaction)
		assert.Equal(t, err, nil)
		payableRepo := repository.NewPayableRepository(DbConn)
		err = payableRepo.Save(context.Background(), firstPayable)
		assert.Equal(t, err, nil)
		err = payableRepo.Save(context.Background(), secondPayable)
		assert.Equal(t, err, nil)
		url := fmt.Sprintf("/balance/%s", clientId)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Equal(t, err, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var output usecase.ClientBalanceOutput
		err = json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Equal(t, err, nil)
		assert.Equal(t, output, usecase.ClientBalanceOutput{
			Balance: usecase.ClientBalanceStatus{
				Paid:      0,
				WaitFunds: 190,
			},
		})
	})

	t.Run("Should get client balance that with payables with status 'available' and 'wait_funds'", func(t *testing.T) {
		defer cleanDatabaseTables(DbConn)
		card, err := entity.NewCard("Teste da Silva", "123", "1111-1111-1111-1111", time.Now().AddDate(5, 0, 0))
		assert.Equal(t, err, nil)
		clientId := uuid.New().String()
		firstTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 153.35, time.Now(), card, entity.PaymentMethods["credit_card"])
		assert.Equal(t, err, nil)
		secondTransaction, err := entity.NewTransaction(uuid.New().String(), clientId, "compra teste", 103.5, time.Now(), card, entity.PaymentMethods["debit_card"])
		assert.Equal(t, err, nil)
		firstPayable, err := entity.PayableFactory(firstTransaction)
		assert.Equal(t, err, nil)
		secondPayable, err := entity.PayableFactory(secondTransaction)
		assert.Equal(t, err, nil)
		payableRepo := repository.NewPayableRepository(DbConn)
		err = payableRepo.Save(context.Background(), firstPayable)
		assert.Equal(t, err, nil)
		err = payableRepo.Save(context.Background(), secondPayable)
		url := fmt.Sprintf("/balance/%s", clientId)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Equal(t, err, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var output usecase.ClientBalanceOutput
		err = json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Equal(t, err, nil)
		assert.Equal(t, usecase.ClientBalanceOutput{
			Balance: usecase.ClientBalanceStatus{
				Paid:      100,
				WaitFunds: 145,
			},
		}, output)
	})
}
