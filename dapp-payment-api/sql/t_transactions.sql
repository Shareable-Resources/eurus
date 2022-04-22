CREATE TABLE IF NOT EXISTS payment.t_transactions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submit_time timestamp with time zone NOT NULL,
    confirmed_time timestamp with time zone NOT NULL,
    network_id bigint NOT NULL,
    token_id bigint NOT NULL,
    from_address text NOT NULL,
    amount numeric(78) NOT NULL,
    merchant_id bigint NOT NULL,
    tag text NOT NULL,
    merchant_seq_no bigint NOT NULL,
    submission_id bigint NOT NULL,
    onchain_status int NOT NULL DEFAULT 0,
    confirm_status int NOT NULL DEFAULT 0,
    signature text NOT NULL,
    tx_hash text NOT NULL,
    block_hash text NOT NULL,
    block_number bigint NOT NULL,
    CONSTRAINT t_transactions_pkey PRIMARY KEY (id),
    CONSTRAINT t_transactions_merchant_id_merchant_seq_no_key UNIQUE (merchant_id, merchant_seq_no),
    CONSTRAINT t_transactions_network_id_tx_hash_key UNIQUE (network_id, tx_hash)
);

CREATE INDEX IF NOT EXISTS t_transactions_confirmed_time_index
    ON payment.t_transactions USING btree
    (confirmed_time ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS t_transactions_submit_time_index
    ON payment.t_transactions USING btree
    (submit_time ASC NULLS LAST);
