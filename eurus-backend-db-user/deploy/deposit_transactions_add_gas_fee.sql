-- Deploy eurus-backend-db-user:deposit_transactions_add_gas_fee to pg

BEGIN;

    ALTER TABLE deposit_transactions ADD COLUMN mainnet_gas_fee numeric(78);

COMMIT;
