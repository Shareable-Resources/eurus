-- Revert eurus-backend-db-user:transaction_index_schema from pg

BEGIN;

DROP TABLE transaction_indexs;

COMMIT;
