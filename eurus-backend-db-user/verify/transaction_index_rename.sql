-- Verify eurus-backend-db-user:transaction_index_rename on pg

BEGIN;

SELECT * FROM transaction_indices LIMIT 1;

ROLLBACK;
