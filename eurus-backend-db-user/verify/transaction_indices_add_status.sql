-- Verify eurus-backend-db-user:transaction_indices_add_status on pg

BEGIN;

SELECT * FROM transaction_indices limit 1;

ROLLBACK;
