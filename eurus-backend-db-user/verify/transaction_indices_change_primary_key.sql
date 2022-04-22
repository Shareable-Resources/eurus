-- Verify eurus-backend-db-user:transaction_indices_change_primary_key on pg

BEGIN;

SELECT * FROM transaction_indices limit 1;

ROLLBACK;
