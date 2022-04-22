-- Revert eurus-backend-db-user:deposit_transactions_drop_innet_trans_hash from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE deposit_transactions ADD COLUMN innet_trans_hash VARCHAR(255);

COMMIT;
