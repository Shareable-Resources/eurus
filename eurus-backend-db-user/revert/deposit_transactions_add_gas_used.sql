-- Revert eurus-backend-db-user:deposit_transactions_add_gas_used from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE deposit_transactions DROP COLUMN mainnet_gas_used;

COMMIT;
