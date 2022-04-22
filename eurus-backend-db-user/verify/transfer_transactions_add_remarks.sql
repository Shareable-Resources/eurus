-- Verify eurus-backend-db-user:transfer_transactions_add_remarks on pg

BEGIN;

-- XXX Add verifications here.
SELECT remarks from transfer_transactions limit 1;

ROLLBACK;
