-- Verify eurus-backend-db-user:transfer_transactions_add_status on pg

BEGIN;

-- XXX Add verifications here.
SELECT status from transfer_transactions LIMIT 1;

ROLLBACK;
