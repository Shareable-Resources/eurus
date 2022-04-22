-- Revert eurus-backend-db-user:deposit_transactions_add_collect_trans from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE deposit_transactions DROP COLUMN mainnet_collect_trans_hash;

COMMIT;
