-- Deploy eurus-backend-db-user:transfer_transactions_add_to_address to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS to_address VARCHAR(255);

COMMIT;
