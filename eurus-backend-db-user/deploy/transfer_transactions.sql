-- Deploy eurus-backend-db-user:transfer_transactions to pg

BEGIN;

CREATE TABLE IF NOT EXISTS transfer_transactions (
    user_id BIGINT NOT NULL,
    asset_name VARCHAR(50) NOT NULL,
    wallet_address VARCHAR(50) NOT NULL,
    tx_hash VARCHAR(50) NOT NULL,
    chain SMALLINT NOT NULL,
    amount numeric(78) NOT NULL,
    gas_fee numeric(78) NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_date TIMESTAMP WITH TIME ZONE NOT NULL
);

 ALTER TABLE transfer_transactions DROP CONSTRAINT IF EXISTS transfer_transactions_pkey;
  ALTER TABLE transfer_transactions ADD CONSTRAINT transfer_transactions_pkey primary key (tx_hash);




COMMIT;
