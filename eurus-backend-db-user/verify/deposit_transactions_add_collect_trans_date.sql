-- Verify eurus-backend-db-user:deposit_transactions_add_collect_trans_date on pg

BEGIN;

-- XXX Add verifications here.
SELECT mainnet_collect_trans_date FROM deposit_transactions LIMIT 1;

ROLLBACK;
