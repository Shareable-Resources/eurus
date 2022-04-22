-- Deploy eurus-backend-db-user:deposit_transactions_drop_innet_trans_hash to pg

BEGIN;

-- XXX Add DDLs here.

ALTER TABLE deposit_transactions DROP COLUMN innet_trans_hash;

COMMIT;
