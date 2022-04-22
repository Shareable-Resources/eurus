-- Revert eurus-backend-db-user:transaction_indices_add_status from pg

BEGIN;

ALTER TABLE transaction_indices DROP COLUMN status;

COMMIT;
