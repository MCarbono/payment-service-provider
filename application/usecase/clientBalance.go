package usecase

import (
	"context"
	"database/sql"
	"fmt"
	db "payment-service-provider/infra/db/sqlc"
	"payment-service-provider/infra/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := tracing.Tracer.Start(ctx, "making query to database")
	span.AddEvent("query", trace.WithAttributes(attribute.String("name", "GetBalanceByStatuses")))
	span.SetAttributes(attribute.String("client_id", input.ClientID))
	balance, err := uc.conn.GetBalanceByStatuses(ctx, sql.NullString{String: input.ClientID, Valid: true})
	if err != nil {
		err = fmt.Errorf("error trying to retrieve balance from user %s: %w", input.ClientID, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, err
	}
	span.End()
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
	Paid      float32 `json:"available"`
	WaitFunds float32 `json:"wait_funds"`
}
