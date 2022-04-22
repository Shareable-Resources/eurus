-- Deploy eurus-backend-db-report:20220214_add_wallet_balance_configs to pg

BEGIN;

CREATE TABLE IF NOT EXISTS wallet_balance_configs (
	id bigserial PRIMARY KEY,
	service_id bigint,
	custom_wallet_address varchar(50),
	custom_wallet_type int,
	chain_id bigint NOT NULL,
	asset_name varchar(100) NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL
);

COMMIT;
