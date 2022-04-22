-- Verify eurus-backend-db-user:transaction_index_schema on pg

BEGIN;

SELECT * FROM transaction_indexs;

ROLLBACK;
