-- Deploy eurus-backend-db-user:transfer_transactions_remove_wrong_idx to pg

BEGIN;

-- XXX Add DDLs here.
DROP INDEX transfer_transactions_request_trans_id_idx;

CREATE INDEX IF NOT EXISTS transfer_transactions_request_trans_id_idx ON transfer_transactions(from_address, request_trans_id, user_id, chain);

COMMIT;
