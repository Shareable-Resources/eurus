-- Verify eurus-backend-db-user:transfer_transactions_add_to_address on pg

BEGIN;

-- XXX Add verifications here.
SELECT to_address FROM transfer_transactions LIMIT 1;

ROLLBACK;
