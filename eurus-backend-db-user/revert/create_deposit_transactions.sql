-- Revert eurus-backend-db-user:create_deposit_transactions from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE deposit_transactions;

COMMIT;
