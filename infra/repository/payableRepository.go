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

func NewPayableRepositoryWithTx(TX *sql.Tx) *PayableRepository {
	queries := db.New(TX).WithTx(TX)
	return &PayableRepository{
		queries: queries,
	}
}

func NewPayableRepository(DB *sql.DB) *PayableRepository {
	queries := db.New(DB)
	return &PayableRepository{
		queries: queries,
	}
}

func (r *PayableRepository) Save(ctx context.Context, payable entity.PayableInterface) error {
	err := r.queries.CreatePayable(ctx, db.CreatePayableParams{
		ID:            payable.GetData().GetID(),
		ClientID:      sql.NullString{String: payable.GetData().GetClientID(), Valid: true},
		TransactionID: sql.NullString{String: payable.GetData().GetTransactionID(), Valid: true},
		Status:        sql.NullString{String: payable.GetData().GetStatus(), Valid: true},
		FeeAmount:     sql.NullInt64{Int64: int64(payable.GetData().GetFeeAmount() * 100), Valid: true},
		Amount:        sql.NullInt64{Int64: int64(payable.GetData().GetAmount() * 100), Valid: true},
		PaymentDate:   payable.GetData().GetPaymentDate(),
		CreatedAt:     payable.GetData().GetCreatedAt(),
	})
	if err != nil {
		return fmt.Errorf("error trying to save payable: %w", err)
	}
	return nil
}

func (r *PayableRepository) GetByID(ctx context.Context, ID string) (entity.PayableInterface, error) {
	model, err := r.queries.GetPayableByID(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPayableNotFound
		}
		return nil, fmt.Errorf("error trying to retrieve payable from DB: %w", err)
	}
	status := entity.GetPayableStatus(model.Status.String)
	payable := entity.RecoverPayableData(
		model.ID,
		model.ClientID.String,
		model.TransactionID.String,
		status,
		float32(model.FeeAmount.Int64)/100,
		float32(model.Amount.Int64)/100,
		model.CreatedAt,
		model.PaymentDate,
	)
	return entity.RestorePayable(payable), nil
}
