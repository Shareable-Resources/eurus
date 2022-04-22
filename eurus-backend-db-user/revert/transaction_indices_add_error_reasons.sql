-- Revert eurus-backend-db-user:transaction_indices_add_error_reasons from pg

BEGIN;

ALTER TABLE transaction_indices DROP COLUMN error_reasons;

COMMIT;
