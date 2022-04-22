-- Verify eurus-backend-db-user:transaction_indices_status_constraint on pg

BEGIN;

SELECT * FROM transaction_indices limit 1;
ROLLBACK;
