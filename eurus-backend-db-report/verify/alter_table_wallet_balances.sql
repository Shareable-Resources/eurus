-- Verify eurus-backend-db-report:alter_table_wallet_balances on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM wallet_balances LIMIT 5;

ROLLBACK;
