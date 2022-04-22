-- Revert eurus-backend-db-user:pending_prewithdraws_add_idx_request_trans_hash_request_trans_id from pg

BEGIN;

    SELECT * FROM pending_prewithdraws;
    DROP INDEX idx_request_trans_hash_request_trans_id;

COMMIT;
