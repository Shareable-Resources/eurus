-- Deploy eurus-backend-db-user:transfer_transactions_add_column to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions RENAME COLUMN wallet_address TO from_address;

ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS confirm_trans_hash VARCHAR(255);

ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS last_modified_date TIMESTAMPTZ NOT NULL DEFAULT ('1970-01-01');

COMMIT;
