-- Deploy eurus-backend-db-user:create_deposit_transactions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE deposit_transactions (
	id SERIAL PRIMARY KEY,
	mainnet_trans_hash VARCHAR(255) NOT NULL,
	amount numeric(78) NOT NULL,
	asset_name VARCHAR(20) NOT NULL,
	mainnet_from_address VARCHAR(255),
	mainnet_to_address VARCHAR(255),
	mainnet_trans_date TIMESTAMPTZ,
	mint_trans_id BIGINT,
	mint_trans_hash VARCHAR(200),
	mint_date TIMESTAMPTZ,
	innet_trans_hash VARCHAR(255),
	innet_to_address VARCHAR(255),
	innet_from_address VARCHAR(255),
	customer_id	BIGINT,
	customer_type SMALLINT, 
	status	SMALLINT,
	remarks TEXT,
	created_date TIMESTAMPTZ NOT NULL,
	last_modified_date TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX idx_deposit_transactions_mainnet_trans_hash ON deposit_transactions(mainnet_trans_hash);

CREATE INDEX idx_deposit_transactions_customer_id_customer_type ON deposit_transactions(customer_id, customer_type);

CREATE INDEX idx_deposit_transactions_created_date ON deposit_transactions(created_date);

COMMIT;
