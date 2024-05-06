CREATE TABLE transactions (
    id varchar PRIMARY KEY,
    client_id varchar,
    description varchar,
    value bigint,
    card_owner_name varchar,
    card_verification_code varchar,
    card_last_digits varchar,
    card_valid_date timestamptz NOT NULL,
    payment_method varchar,
    created_at timestamptz NOT NULL
);

