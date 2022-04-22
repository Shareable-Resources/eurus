-- Revert eurus-backend-db-user:withdraw_transactions_add_remarks from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transactions DROP COLUMN remarks;

COMMIT;
