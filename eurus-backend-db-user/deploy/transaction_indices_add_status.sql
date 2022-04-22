-- Deploy eurus-backend-db-user:transaction_indices_add_status to pg

BEGIN;

ALTER TABLE transaction_indices ADD COLUMN IF NOT EXISTS status BOOLEAN;
COMMIT;
