-- Revert eurus-backend-db-user:withdraw_requests from pg

BEGIN;

DROP TABLE withdraw_requests;
DROP INDEX idx_request_trans_hash_service_id;

COMMIT;
