-- Verify eurus-backend-db-user:withdraw_transaction on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM withdraw_transactions LIMIT 1;


ROLLBACK;
