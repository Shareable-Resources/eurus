-- Deploy eurus-backend-db-user:transfer_transactions_add_request_trans_id to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS request_trans_id bigint;
ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS is_send boolean NOT NULL DEFAULT(true);

ALTER TABLE transfer_transactions DROP CONSTRAINT transfer_transactions_pkey;
ALTER TABLE transfer_transactions ADD PRIMARY KEY (user_id, tx_hash, chain);
CREATE UNIQUE INDEX transfer_transactions_request_trans_id_idx ON transfer_transactions(wallet_address, request_trans_id, user_id, chain);
CREATE INDEX transfer_transactions_request_trans_id_idx1 ON transfer_transactions(wallet_address, request_trans_id, chain);


COMMIT;
