-- Revert eurus-backend-db-user:transfer_transactions_remove_wrong_idx from pg

BEGIN;

-- XXX Add DDLs here.
CREATE INDEX IF NOT EXISTS transfer_transactions_request_trans_id_idx ON transfer_transactions(from_address, request_trans_id, user_id, chain);


COMMIT;
