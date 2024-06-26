// Code generated by sqlc. DO NOT EDIT.
// source: payables.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPayable = `-- name: CreatePayable :exec
INSERT INTO payables (id, client_id, transaction_id, status, fee_amount, amount, payment_date, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreatePayableParams struct {
	ID            string
	ClientID      sql.NullString
	TransactionID sql.NullString
	Status        sql.NullString
	FeeAmount     sql.NullInt64
	Amount        sql.NullInt64
	PaymentDate   time.Time
	CreatedAt     time.Time
}

func (q *Queries) CreatePayable(ctx context.Context, arg CreatePayableParams) error {
	_, err := q.db.ExecContext(ctx, createPayable,
		arg.ID,
		arg.ClientID,
		arg.TransactionID,
		arg.Status,
		arg.FeeAmount,
		arg.Amount,
		arg.PaymentDate,
		arg.CreatedAt,
	)
	return err
}

const getBalanceByStatuses = `-- name: GetBalanceByStatuses :many
SELECT s.status::text as status, COALESCE(SUM(t.amount), 0)::bigint AS total_amount
FROM (
    SELECT 'paid' AS status
    UNION
    SELECT 'wait_funds' AS status
) s
LEFT JOIN payables t ON s.status = t.status AND t.client_id = $1
GROUP BY s.status
`

type GetBalanceByStatusesRow struct {
	Status      string
	TotalAmount int64
}

func (q *Queries) GetBalanceByStatuses(ctx context.Context, clientID sql.NullString) ([]GetBalanceByStatusesRow, error) {
	rows, err := q.db.QueryContext(ctx, getBalanceByStatuses, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBalanceByStatusesRow{}
	for rows.Next() {
		var i GetBalanceByStatusesRow
		if err := rows.Scan(&i.Status, &i.TotalAmount); err != nil {
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

const getPayableByID = `-- name: GetPayableByID :one
SELECT id, client_id, transaction_id, status, fee_amount, amount, payment_date, created_at FROM payables
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPayableByID(ctx context.Context, id string) (Payable, error) {
	row := q.db.QueryRowContext(ctx, getPayableByID, id)
	var i Payable
	err := row.Scan(
		&i.ID,
		&i.ClientID,
		&i.TransactionID,
		&i.Status,
		&i.FeeAmount,
		&i.Amount,
		&i.PaymentDate,
		&i.CreatedAt,
	)
	return i, err
}
