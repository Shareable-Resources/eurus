-- Deploy eurus-backend-db-user:pending_prewithdraws_delete_idx_approval_wallet_address_request_trans_id to pg

BEGIN;

    SELECT * FROM pending_prewithdraws;
    DROP INDEX idx_approval_wallet_address_request_trans_id;
    
COMMIT;
