-- Revert eurus-backend-db-user:withdraw_transaction_rename from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transactions RENAME TO withdraw_transaction;

COMMIT;
