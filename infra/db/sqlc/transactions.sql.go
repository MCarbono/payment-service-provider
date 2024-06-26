// Code generated by sqlc. DO NOT EDIT.
// source: transactions.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (
    id, 
    client_id,
    description,
    value,
    card_owner_name,
    card_verification_code,
    card_last_digits,
    card_valid_date,
    payment_method,
    created_at
    ) VALUES(
     $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
`

type CreateTransactionParams struct {
	ID                   string
	ClientID             sql.NullString
	Description          sql.NullString
	Value                sql.NullInt64
	CardOwnerName        sql.NullString
	CardVerificationCode sql.NullString
	CardLastDigits       sql.NullString
	CardValidDate        time.Time
	PaymentMethod        sql.NullString
	CreatedAt            time.Time
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, createTransaction,
		arg.ID,
		arg.ClientID,
		arg.Description,
		arg.Value,
		arg.CardOwnerName,
		arg.CardVerificationCode,
		arg.CardLastDigits,
		arg.CardValidDate,
		arg.PaymentMethod,
		arg.CreatedAt,
	)
	return err
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT id, client_id, description, value, card_owner_name, card_verification_code, card_last_digits, card_valid_date, payment_method, created_at
FROM transactions WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransactionByID(ctx context.Context, id string) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransactionByID, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.Description,
		&i.Value,
		&i.CardOwnerName,
		&i.CardVerificationCode,
		&i.CardLastDigits,
		&i.CardValidDate,
		&i.PaymentMethod,
		&i.CreatedAt,
	)
	return i, err
}

const getTransactionsByClientID = `-- name: GetTransactionsByClientID :many
SELECT id, client_id, description, value, card_owner_name, card_verification_code, card_last_digits, card_valid_date, payment_method, created_at
FROM transactions WHERE client_id = $1
`

func (q *Queries) GetTransactionsByClientID(ctx context.Context, clientID sql.NullString) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByClientID, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.ClientID,
			&i.Description,
			&i.Value,
			&i.CardOwnerName,
			&i.CardVerificationCode,
			&i.CardLastDigits,
			&i.CardValidDate,
			&i.PaymentMethod,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
