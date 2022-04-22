-- Revert eurus-backend-db-user:transaction_indices_status_constraint from pg

BEGIN;

ALTER TABLE transaction_indices ALTER COLUMN status SET NULL;

COMMIT;
