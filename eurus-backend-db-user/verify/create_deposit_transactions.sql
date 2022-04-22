-- Verify eurus-backend-db-user:create_deposit_transactions on pg

BEGIN;

-- XXX Add verifications here.
select * from deposit_transactions LIMIT 1;

ROLLBACK;
