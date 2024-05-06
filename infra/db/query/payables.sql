-- name: CreatePayable :exec
INSERT INTO payables (id, client_id, transaction_id, status, fee_amount, amount, payment_date, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetPayableByID :one
SELECT id, client_id, transaction_id, status, fee_amount, amount, payment_date, created_at FROM payables
WHERE id = $1 LIMIT 1;

-- name: GetBalanceByStatuses :many
SELECT s.status::text as status, COALESCE(SUM(t.amount), 0)::bigint AS total_amount
FROM (
    SELECT 'paid' AS status
    UNION
    SELECT 'wait_funds' AS status
) s
LEFT JOIN payables t ON s.status = t.status AND t.client_id = $1
GROUP BY s.status;