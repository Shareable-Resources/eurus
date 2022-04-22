-- Deploy eurus-backend-db-user:deposit_transactions_add_gas_used to pg

BEGIN;

-- XXX Add DDLs here.
	ALTER TABLE deposit_transactions ADD COLUMN mainnet_gas_used numeric(78);

COMMIT;
