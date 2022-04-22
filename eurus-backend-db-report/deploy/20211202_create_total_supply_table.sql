-- Deploy eurus-backend-db-report:20211202_create_total_supply_table to pg

BEGIN;

CREATE TABLE IF NOT EXISTS asset_total_supplies (
	id SERIAL PRIMARY KEY,
	asset_name VARCHAR(100) NOT NULL,
	total_supply NUMERIC(78),
	chain_id int NOT NULL,
	block_number BIGINT,
	asset_address VARCHAR(50),
	mark_date timestamptz NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS asset_total_supplies_idx1 ON asset_total_supplies(asset_name, chain_id, mark_date);
CREATE INDEX IF NOT EXISTS asset_total_supplies_idx2 ON asset_total_supplies (mark_date);


COMMIT;
