CREATE TABLE IF NOT EXISTS payment.t_networks
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    network_code character varying(32) NOT NULL,
    network_name character varying(128) NOT NULL,
    chain_id integer,
    rpc_url text,
    CONSTRAINT t_networks_pkey PRIMARY KEY (id),
    CONSTRAINT t_networks_network_code_key UNIQUE (network_code)
)
