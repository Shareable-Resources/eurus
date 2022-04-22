-- Verify eurus-backend-db-user:deposit_transactions_add_gas_used on pg

BEGIN;

-- XXX Add verifications here.
SELECT mainnet_gas_used from deposit_transactions LIMIT 1;

ROLLBACK;
