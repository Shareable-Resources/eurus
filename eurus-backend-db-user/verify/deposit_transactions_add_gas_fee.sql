-- Verify eurus-backend-db-user:deposit_transactions_add_gas_fee on pg

BEGIN;

SELECT mainnet_gas_fee from deposit_transactions LIMIT 1;

ROLLBACK;
