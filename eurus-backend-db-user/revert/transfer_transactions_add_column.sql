-- Revert eurus-backend-db-user:transfer_transactions_add_column from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions RENAME COLUMN from_address TO wallet_address;
COMMIT;
