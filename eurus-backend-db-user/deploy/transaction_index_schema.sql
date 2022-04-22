-- Deploy eurus-backend-db-user:transaction_index_schema to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE transaction_indexs (
	user_id BIGINT PRIMARY KEY,
    wallet_address VARCHAR(100) NOT NULL,
    tx_hash VARCHAR(100) NOT NULL,
    created_date TIMESTAMP WITH TIME ZONE NOT NULL,
    asset_name VARCHAR(100),
    status BOOLEAN NOT NULL
);

COMMIT;
