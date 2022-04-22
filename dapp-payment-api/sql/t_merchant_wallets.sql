CREATE TABLE IF NOT EXISTS payment.t_merchant_wallets
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merchant_id bigint NOT NULL,
    token_id bigint NOT NULL,
    address text NOT NULL,
    CONSTRAINT t_merchant_wallets_pkey PRIMARY KEY (id),
    CONSTRAINT t_merchant_wallets_merchant_id_token_id_address_key UNIQUE (merchant_id, token_id, address)
);

CREATE INDEX IF NOT EXISTS t_merchant_wallets_address_index
    ON payment.t_merchant_wallets USING btree
    (address ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS t_merchant_wallets_merchant_id_token_id_index
    ON payment.t_merchant_wallets USING btree
    (merchant_id ASC NULLS LAST, token_id ASC NULLS LAST);

CREATE INDEX IF NOT EXISTS t_merchant_wallets_token_id_address_index
    ON payment.t_merchant_wallets USING btree
    (token_id ASC NULLS LAST, address ASC NULLS LAST);
