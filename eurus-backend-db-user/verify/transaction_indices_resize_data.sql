-- Verify eurus-backend-db-user:transaction_indices_resize_data on pg

BEGIN;

SELECT * FROM transaction_indices limit 1;

ROLLBACK;
