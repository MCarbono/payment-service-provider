package usecase

import (
	"context"
	"database/sql"
	"fmt"
	db "payment-service-provider/infra/db/sqlc"
)

type ClientBalance struct {
	conn db.Querier
}

func NewClientBalance(db db.Querier) *ClientBalance {
	return &ClientBalance{
		conn: db,
	}
}

func (uc *ClientBalance) Execute(ctx context.Context, input *ClientBalanceInput) (*ClientBalanceOutput, error) {
	balance, err := uc.conn.GetBalanceByStatuses(ctx, sql.NullString{String: input.ClientID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("error trying to retrieve balance from user %s: %w", input.ClientID, err)
	}

	return &ClientBalanceOutput{
		Balance: ClientBalanceStatus{
			Paid:      float32(balance[0].TotalAmount / 100),
			WaitFunds: float32(balance[1].TotalAmount / 100),
		},
	}, nil
}

type ClientBalanceInput struct {
	ClientID string `json:"client_id"`
}

type ClientBalanceOutput struct {
	Balance ClientBalanceStatus `json:"balance"`
}

type ClientBalanceStatus struct {
	Paid      float32 `json:"paid"`
	WaitFunds float32 `json:"wait_funds"`
}
