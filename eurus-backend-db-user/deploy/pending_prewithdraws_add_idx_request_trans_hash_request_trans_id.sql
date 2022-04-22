-- Deploy eurus-backend-db-user:pending_prewithdraws_add_idx_request_trans_hash_request_trans_id to pg

BEGIN;

CREATE UNIQUE INDEX idx_request_trans_hash_request_trans_id ON pending_prewithdraws(request_trans_hash, request_trans_id);

COMMIT;
