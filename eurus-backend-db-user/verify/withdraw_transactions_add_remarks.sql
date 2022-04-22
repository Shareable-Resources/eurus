-- Verify eurus-backend-db-user:withdraw_transactions_add_remarks on pg

BEGIN;

-- XXX Add verifications here.
SELECT remarks FROM withdraw_transactions LIMIT 1;

ROLLBACK;
