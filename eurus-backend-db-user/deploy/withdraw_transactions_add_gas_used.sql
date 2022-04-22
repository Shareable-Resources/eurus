-- Deploy eurus-backend-db-user:withdraw_transactions_add_gas_used to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transactions ADD COLUMN sidechain_gas_used numeric(78);

COMMIT;
