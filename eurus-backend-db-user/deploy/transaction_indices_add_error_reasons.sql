-- Deploy eurus-backend-db-user:transaction_indices_add_error_reasons to pg

BEGIN;

ALTER TABLE transaction_indices ADD COLUMN error_reasons VARCHAR(200);

COMMIT;
