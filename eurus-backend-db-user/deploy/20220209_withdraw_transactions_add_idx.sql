-- Deploy eurus-backend-db-user:20220209_withdraw_transactions_add_idx to pg

BEGIN;

-- XXX Add DDLs here.
CREATE UNIQUE INDEX IF NOT EXISTS  withdraw_transactions_idx1 ON withdraw_transactions (request_trans_hash);

COMMIT;
