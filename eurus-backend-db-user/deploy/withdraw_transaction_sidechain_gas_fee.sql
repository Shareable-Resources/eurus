-- Deploy eurus-backend-db-user:withdraw_transaction_sidechain_gas_fee to pg

BEGIN;

ALTER TABLE withdraw_transactions ADD COLUMN sidechain_gas_fee numeric(78);

COMMIT;
