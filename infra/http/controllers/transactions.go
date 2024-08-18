package controllers

import (
	"encoding/json"
	"net/http"
	"payment-service-provider/application/usecase"
	"payment-service-provider/infra/tracing"

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
	ctx, span := tracing.Tracer.Start(r.Context(), "controller - processing transaction")
	defer span.End()
	var input usecase.ProcessTransactionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()
	_, err = c.processTransaction.Execute(ctx, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newControllerOutput(ctx, nil, err))
		return
	}
	output := newControllerOutput(ctx, nil, nil)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&output)
}

func (c *Transaction) ListTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.Tracer.Start(r.Context(), "controller - listing all transactions")
	defer span.End()
	id := chi.URLParam(r, "client_id")
	output, err := c.listTransactions.Execute(ctx, &usecase.ListTransationsInput{ClientID: id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newControllerOutput(ctx, nil, err))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newControllerOutput(ctx, output, err))
}

func (c *Transaction) Balance(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.Tracer.Start(r.Context(), "controller - client balance")
	defer span.End()
	id := chi.URLParam(r, "client_id")
	output, err := c.clientBalance.Execute(ctx, &usecase.ClientBalanceInput{ClientID: id})
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
