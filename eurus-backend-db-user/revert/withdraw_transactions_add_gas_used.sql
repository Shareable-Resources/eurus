-- Revert eurus-backend-db-user:withdraw_transactions_add_gas_used from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transactions DROP COLUMN sidechain_gas_used;


COMMIT;
