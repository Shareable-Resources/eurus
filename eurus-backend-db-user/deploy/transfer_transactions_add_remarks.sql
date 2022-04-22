-- Deploy eurus-backend-db-user:transfer_transactions_add_remarks to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN IF NOT EXISTS remarks TEXT NULL;

COMMIT;
