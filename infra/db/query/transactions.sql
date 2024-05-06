-- name: CreateTransaction :exec
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
);

-- name: GetTransactionByID :one
SELECT id, client_id, description, value, card_owner_name, card_verification_code, card_last_digits, card_valid_date, payment_method, created_at
FROM transactions WHERE id = $1 LIMIT 1;

-- name: GetTransactionsByClientID :many
SELECT id, client_id, description, value, card_owner_name, card_verification_code, card_last_digits, card_valid_date, payment_method, created_at
FROM transactions WHERE client_id = $1;