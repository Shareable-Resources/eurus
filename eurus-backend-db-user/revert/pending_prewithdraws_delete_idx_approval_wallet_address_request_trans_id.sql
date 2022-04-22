-- Revert eurus-backend-db-user:pending_prewithdraws_delete_idx_approval_wallet_address_request_trans_id from pg

BEGIN;

    CREATE UNIQUE INDEX idx_approval_wallet_address_request_trans_id ON pending_prewithdraws(approval_wallet_address, request_trans_id);

COMMIT;
