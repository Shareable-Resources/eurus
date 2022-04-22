-- Deploy eurus-backend-db-user:deposit_transactions_add_collect_trans to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE deposit_transactions ADD COLUMN mainnet_collect_trans_hash VARCHAR(255);

COMMIT;
