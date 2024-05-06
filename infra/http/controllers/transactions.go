package controllers

import (
	"encoding/json"
	"net/http"
	"payment-service-provider/application/usecase"
)

type TransactionController struct {
	processTransaction *usecase.ProcessTransaction
	listTransactions   *usecase.ListTransactions
	clientBalance      *usecase.ClientBalance
}

func NewTransactionController(
	processTransaction *usecase.ProcessTransaction,
	listTransactions *usecase.ListTransactions,
	clientBalance *usecase.ClientBalance,
) *TransactionController {
	return &TransactionController{
		processTransaction: processTransaction,
		listTransactions:   listTransactions,
		clientBalance:      clientBalance,
	}
}

func (c *TransactionController) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	var input usecase.ProcessTransactionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
}

func (c *TransactionController) ListTransactions(w http.ResponseWriter, r *http.Request) {

}

func (c *TransactionController) Balance(w http.ResponseWriter, r *http.Request) {

}
