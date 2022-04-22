-- Deploy eurus-backend-db-user:20211102_topup_transactions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS top_up_transactions (
	tx_hash varchar(255) PRIMARY KEY,
	customer_id bigint NOT NULL,
	customer_type smallint NOT NULL,
	from_address varchar(255) NOT NULL,
	to_address varchar(255) NOT NULL,
	transfer_gas numeric(78) NOT NULL,
	target_gas numeric(78) NOT NULL,
	status smallint NOT NULL,
	is_direct_top_up bool NOT NULL,
	remarks text,
	trans_gas_used bigint,
	user_gas_used bigint,
	gas_price numeric(78),
	transaction_date timestamptz NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL
);

CREATE INDEX top_up_transactions_idx1 ON top_up_transactions (customer_id, customer_type);
CREATE INDEX top_up_transactions_idx2 ON top_up_transactions (from_address);
CREATE INDEX top_up_transactions_idx3 ON top_up_transactions (created_date);


COMMIT;
