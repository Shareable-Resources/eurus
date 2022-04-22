-- Verify eurus-backend-db-user:20210903_create_unique_index on pg

BEGIN;

-- XXX Add verifications here.
SELECT * from users LIMIT 1;
SELECT * from deposit_transactions LIMIT 1;
SELECT * from withdraw_transactions LIMIT 1;
SELECT * from transfer_transactions LIMIT 1;
SELECT * from user_faucets LIMIT 1;
SELECT * from verifications LIMIT 1;
SELECT * from user_faucets LIMIT 1;

ROLLBACK;
