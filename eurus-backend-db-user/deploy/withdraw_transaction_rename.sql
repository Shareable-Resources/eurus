-- Deploy eurus-backend-db-user:withdraw_transaction_rename to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transaction RENAME TO withdraw_transactions;
COMMIT;
