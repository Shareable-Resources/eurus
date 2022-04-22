-- Verify eurus-backend-db-user:transaction_index_change_fields on pg

BEGIN;

SELECT * FROM transaction_indices limit 1;

ROLLBACK;
