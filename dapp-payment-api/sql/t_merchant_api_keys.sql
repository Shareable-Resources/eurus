CREATE TABLE IF NOT EXISTS payment.t_merchant_api_keys
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merchant_id bigint NOT NULL,
    api_key text NOT NULL,
	salt text NOT NULL,
    CONSTRAINT t_merchant_api_keys_pkey PRIMARY KEY (id),
    CONSTRAINT t_merchant_api_keys_merchant_id_api_key_key UNIQUE (merchant_id, api_key)
);

CREATE INDEX IF NOT EXISTS t_merchant_api_keys_merchant_id_index
    ON payment.t_merchant_api_keys USING btree
    (merchant_id ASC NULLS LAST);