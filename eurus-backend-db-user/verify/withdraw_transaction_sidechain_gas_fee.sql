-- Verify eurus-backend-db-user:withdraw_transaction_sidechain_gas_fee on pg

BEGIN;

SELECT sidechain_gas_fee from withdraw_transactions LIMIT 1;

ROLLBACK;
