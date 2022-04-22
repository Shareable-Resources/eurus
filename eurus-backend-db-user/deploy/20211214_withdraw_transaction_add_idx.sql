-- Deploy eurus-backend-db-user:20211214_withdraw_transaction_add_idx to pg

BEGIN;

-- XXX Add DDLs here.
CREATE INDEX IF NOT EXISTS withdraw_transactions_idx2 ON withdraw_transactions (approval_wallet_address, request_trans_id);

COMMIT;
