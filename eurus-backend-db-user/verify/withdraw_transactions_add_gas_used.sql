-- Verify eurus-backend-db-user:withdraw_transactions_add_gas_used on pg

BEGIN;

-- XXX Add verifications here.
SELECT sidechain_gas_used from withdraw_transactions LIMIT 1;

ROLLBACK;
