-- Revert eurus-backend-db-user:withdraw_transaction_sidechain_gas_fee from pg

BEGIN;

    ALTER TABLE withdraw_transactions DROP COLUMN sidechain_gas_fee;

COMMIT;
