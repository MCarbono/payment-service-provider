CREATE TABLE payables (
    id varchar PRIMARY KEY,
    client_id varchar,
    transaction_id varchar,
    status varchar,
    fee_amount bigint,
    amount bigint,
    payment_date timestamptz NOT NULL,
    created_at timestamptz NOT NULL
);