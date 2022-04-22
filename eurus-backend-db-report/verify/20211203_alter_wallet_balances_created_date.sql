-- Verify eurus-backend-db-report:20211203_alter_wallet_balances_created_date on pg

BEGIN;

-- XXX Add verifications here.
select mark_date from wallet_balances LIMIT 1;

ROLLBACK;
