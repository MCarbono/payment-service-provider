package controllers

import (
	"encoding/json"
	"net/http"
	"payment-service-provider/application/usecase"

	"github.com/go-chi/chi/v5"
)

type Transaction struct {
	processTransaction *usecase.ProcessTransaction
	listTransactions   *usecase.ListTransactions
	clientBalance      *usecase.ClientBalance
}

func NewTransaction(
	processTransaction *usecase.ProcessTransaction,
	listTransactions *usecase.ListTransactions,
	clientBalance *usecase.ClientBalance,
) *Transaction {
	return &Transaction{
		processTransaction: processTransaction,
		listTransactions:   listTransactions,
		clientBalance:      clientBalance,
	}
}

func (c *Transaction) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	var input usecase.ProcessTransactionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()
	_, err = c.processTransaction.Execute(r.Context(), &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonErr{Err: err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *Transaction) ListTransactions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	output, err := c.listTransactions.Execute(r.Context(), &usecase.ListTransationsInput{ClientID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonErr{Err: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&output)
}

func (c *Transaction) Balance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	output, err := c.clientBalance.Execute(r.Context(), &usecase.ClientBalanceInput{ClientID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonErr{Err: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&output)
}

type jsonErr struct {
	Err string `json:"error"`
}
