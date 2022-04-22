CREATE TABLE IF NOT EXISTS payment.t_merchants
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merchant_code character varying(32) NOT NULL,
    merchant_name character varying(128) NOT NULL,
    tag_display_name text,
    tag_description character varying(128),
    merchant_last_seq bigint NOT NULL DEFAULT 0,
    CONSTRAINT t_merchants_pkey PRIMARY KEY (id),
    CONSTRAINT t_merchants_merchant_code_key UNIQUE (merchant_code)
)
