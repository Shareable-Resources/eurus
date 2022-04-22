-- Verify eurus-backend-db-user:withdraw_transaction_rename on pg

BEGIN;

-- XXX Add verifications here.
Select * from withdraw_transactions limit 1;

ROLLBACK;
