-- Deploy eurus-backend-db-user:20210810_create_distributed_tokens to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS distributed_tokens (
	id BIGSERIAL PRIMARY KEY,
	asset_name VARCHAR(100) NOT NULL,
	amount numeric(78) NOT NULL,
	chain smallint,
	distributed_type int NOT NULL,
	user_id bigint,
	tx_hash varchar(255),
	from_address varchar(255) NOT NULL,
	to_address varchar(255) NOT NULL,
	gas_price numeric(78),
	gas_used bigint,
	gas_fee numeric(78),
	created_date timestamp with time zone NOT NULL,
	last_modified_date timestamp with time zone NOT NULL
);

COMMIT;
