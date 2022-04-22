-- Verify eurus-backend-db-report:2021_08_23_alter_wallet_ballance_PK on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM wallet_balances LIMIT 1;

ROLLBACK;
