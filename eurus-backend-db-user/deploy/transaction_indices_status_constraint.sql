-- Deploy eurus-backend-db-user:transaction_indices_status_constraint to pg

BEGIN;

ALTER TABLE transaction_indices ALTER COLUMN status SET NOT NULL;
COMMIT;
