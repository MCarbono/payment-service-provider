package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment-service-provider/domain/entity"
	db "payment-service-provider/infra/db/sqlc"
)

var (
	ErrPayableNotFound = errors.New("payable not found")
)

type PayableRepository struct {
	queries *db.Queries
}

func NewPayableRepository(DB *sql.DB) *PayableRepository {
	queries := db.New(DB)
	return &PayableRepository{
		queries: queries,
	}
}

func (r *PayableRepository) Save(ctx context.Context, payable *entity.Payable) error {
	err := r.queries.CreatePayable(ctx, db.CreatePayableParams{
		ID:            payable.GetID(),
		ClientID:      sql.NullString{String: payable.GetClientID(), Valid: true},
		TransactionID: sql.NullString{String: payable.GetTransactionID(), Valid: true},
		Status:        sql.NullString{String: payable.GetStatus(), Valid: true},
		FeeAmount:     sql.NullInt64{Int64: int64(payable.GetFeeAmount() * 100), Valid: true},
		Amount:        sql.NullInt64{Int64: int64(payable.GetAmount() * 100), Valid: true},
		PaymentDate:   payable.GetPaymentDate(),
		CreatedAt:     payable.GetCreatedAt(),
	})
	if err != nil {
		return fmt.Errorf("error trying to save payable: %w", err)
	}
	return nil
}

func (r *PayableRepository) GetByID(ctx context.Context, ID string) (*entity.Payable, error) {
	model, err := r.queries.GetPayableByID(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPayableNotFound
		}
		return nil, fmt.Errorf("error trying to retrieve payable from DB: %w", err)
	}
	status := entity.GetPayableStatus(model.Status.String)
	payable := entity.NewPayable(
		model.ID,
		model.ClientID.String,
		model.TransactionID.String,
		status,
		float32(model.FeeAmount.Int64)/100,
		float32(model.Amount.Int64)/100,
		model.CreatedAt,
		model.PaymentDate,
	)
	return payable, nil
}
