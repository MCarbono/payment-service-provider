package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment-service-provider/domain/entity"
	db "payment-service-provider/infra/db/sqlc"
	"payment-service-provider/infra/tracing"

	"go.opentelemetry.io/otel/codes"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionRepository struct {
	queries *db.Queries
}

func NewTransactionRepositoryWithTX(TX *sql.Tx) *TransactionRepository {
	queries := db.New(TX).WithTx(TX)
	return &TransactionRepository{
		queries: queries,
	}
}

func NewTransactionRepository(DB *sql.DB) *TransactionRepository {
	queries := db.New(DB)
	return &TransactionRepository{
		queries: queries,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, transaction *entity.Transaction) error {
	_, span := tracing.Tracer.Start(ctx, "saving transaction into database")
	defer span.End()
	card := transaction.GetCard()
	err := r.queries.CreateTransaction(ctx, db.CreateTransactionParams{
		ID:                   transaction.GetID(),
		ClientID:             sql.NullString{String: transaction.GetClientID(), Valid: true},
		Description:          sql.NullString{String: transaction.GetDescription(), Valid: true},
		Value:                sql.NullInt64{Int64: int64(transaction.GetValue()) * 100, Valid: true},
		CardOwnerName:        sql.NullString{String: card.GetOwnerName(), Valid: true},
		CardVerificationCode: sql.NullString{String: card.GetVerificationCode(), Valid: true},
		CardLastDigits:       sql.NullString{String: card.GetLastDigits(), Valid: true},
		CardValidDate:        card.GetValidDate(),
		PaymentMethod:        sql.NullString{String: transaction.PaymentMethod().String(), Valid: true},
		CreatedAt:            transaction.GetCreatedAt(),
	})
	if err != nil {
		err = fmt.Errorf("error trying to insert transaction into the db: %w", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
func (r *TransactionRepository) GetByID(ctx context.Context, ID string) (*entity.Transaction, error) {
	model, err := r.queries.GetTransactionByID(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("error trying to get transation by ID: %w", err)
	}
	paymentMethod := entity.PaymentMethods[model.PaymentMethod.String]
	card := entity.RestoreCard(model.CardOwnerName.String, model.CardVerificationCode.String, model.CardLastDigits.String, model.CardValidDate)
	transaction, err := entity.NewTransaction(
		model.ID,
		model.ClientID.String,
		model.Description.String,
		float32(model.Value.Int64)/100,
		model.CreatedAt,
		card,
		paymentMethod,
	)
	if err != nil {
		return nil, fmt.Errorf("error trying to get transation by ID: %w", err)
	}
	return transaction, nil
}
