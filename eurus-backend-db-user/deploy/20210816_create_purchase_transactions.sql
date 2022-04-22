-- Deploy eurus-backend-db-user:20210816_create_purchase_transactions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS purchase_transactions (
	user_id BIGINT NOT NULL,
	tx_hash VARCHAR(255) NOT NULL,
	chain SMALLINT NOT NULL,
	from_address VARCHAR(255) NOT NULL,
	to_address VARCHAR(255) NOT NULL,
	asset_name VARCHAR(50) NOT NULL,
	amount NUMERIC(78) NOT NULL,
	product_id BIGINT,
	quantity NUMERIC(78),
	gas_fee NUMERIC(78) NOT NULL,
	trans_gas_used BIGINT NOT NULL,
	user_gas_used BIGINT NOT NULL,
	gas_price NUMERIC(78) NOT NULL,
	status SMALLINT NOT NULL,
	purchase_type INT NOT NULL,
	remarks TEXT,
	created_date TIMESTAMPTZ NOT NULL,
	last_modified_date TIMESTAMPTZ NOT NULL,
	PRIMARY KEY (user_id, tx_hash, chain)
);


COMMIT;
