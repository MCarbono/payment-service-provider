package usecase

import (
	"context"
	"database/sql"
	"fmt"
	db "payment-service-provider/infra/db/sqlc"
	"time"
)

type ListTransactions struct {
	conn db.Querier
}

func NewListTransactions(db db.Querier) *ListTransactions {
	return &ListTransactions{
		conn: db,
	}
}

func (uc *ListTransactions) Execute(ctx context.Context, input *ListTransationsInput) ([]ListTransactionsOutput, error) {
	models, err := uc.conn.GetTransactionsByClientID(ctx, sql.NullString{String: input.ClientID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("error trying to get transactions by clientId %s: %w", input.ClientID, err)
	}
	return NewListTransactionsOutput(models), nil
}

type ListTransationsInput struct {
	ClientID string `json:"client_id"`
}

type ListTransactionsOutput struct {
	ID                   string    `json:"id"`
	ClientID             string    `json:"client_id"`
	Description          string    `json:"description"`
	Value                float32   `json:"value"`
	CardOwnerName        string    `json:"card_owner_name"`
	CardVerificationCode string    `json:"card_verification_name"`
	CardLastDigits       string    `json:"card_last_digits"`
	CardValidDate        time.Time `json:"card_valid_date"`
	PaymentMethod        string    `json:"payment_method"`
	CreatedAt            time.Time `json:"created_at"`
}

func NewListTransactionsOutput(transactions []db.Transaction) []ListTransactionsOutput {
	output := make([]ListTransactionsOutput, len(transactions))
	for i := 0; i < len(output); i++ {
		output[i] = ListTransactionsOutput{
			ID:                   transactions[i].ID,
			ClientID:             transactions[i].ClientID.String,
			Description:          transactions[i].Description.String,
			Value:                float32(transactions[i].Value.Int64) / 100,
			CardOwnerName:        transactions[i].CardOwnerName.String,
			CardVerificationCode: transactions[i].CardVerificationCode.String,
			CardLastDigits:       transactions[i].CardLastDigits.String,
			CardValidDate:        transactions[i].CardValidDate,
			PaymentMethod:        transactions[i].PaymentMethod.String,
			CreatedAt:            transactions[i].CreatedAt,
		}
	}
	return output
}
