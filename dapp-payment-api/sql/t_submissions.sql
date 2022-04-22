CREATE TABLE IF NOT EXISTS payment.t_submissions
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submit_time timestamp with time zone NOT NULL,
    network_id bigint NOT NULL,
    token_id bigint NOT NULL,
    from_address text NOT NULL,
    amount numeric(78) NOT NULL,
    merchant_id bigint NOT NULL,
    tag text NOT NULL,
    tx_hash text NOT NULL,
    tx_status int NOT NULL DEFAULT -1,
    payment_status int NOT NULL DEFAULT 0,
    signature text NOT NULL,
    message_body text NOT NULL,
    CONSTRAINT t_submissions_pkey PRIMARY KEY (id),
    CONSTRAINT t_submissions_network_id_tx_hash_signature_key UNIQUE (network_id, tx_hash, signature)
);

CREATE INDEX IF NOT EXISTS t_submissions_submit_time_index
    ON payment.t_submissions USING btree
    (submit_time ASC NULLS LAST);
