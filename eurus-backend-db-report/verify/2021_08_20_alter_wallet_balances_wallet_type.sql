-- Verify eurus-backend-db-report:2021_08_20_alter_wallet_balances_wallet_type on pg

BEGIN;

-- XXX Add verifications here.
SELECT wallet_type from wallet_balances LIMIT 1;

ROLLBACK;
