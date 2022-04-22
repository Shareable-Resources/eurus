CREATE TABLE IF NOT EXISTS payment.t_tokens
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    network_id bigint NOT NULL,
    address text NOT NULL,
    symbol character varying(16) NOT NULL,
    name character varying(64) NOT NULL,
    decimals integer NOT NULL,
    CONSTRAINT t_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT t_tokens_network_id_address_key UNIQUE (network_id, address),
    CONSTRAINT t_tokens_network_id_symbol_key UNIQUE (network_id, symbol)
)
