-- Revert eurus-backend-db-user:deposit_transactions_add_gas_fee from pg

BEGIN;

ALTER TABLE deposit_transactions DROP COLUMN mainnet_gas_fee;

COMMIT;
