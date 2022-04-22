-- Verify eurus-backend-db-user:deposit_transactions_add_collect_trans on pg

BEGIN;

-- XXX Add verifications here.
SELECT mainnet_collect_trans_hash from deposit_transactions LIMIT 1;
ROLLBACK;
