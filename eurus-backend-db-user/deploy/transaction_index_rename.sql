-- Deploy eurus-backend-db-user:transaction_index_rename to pg

BEGIN;

ALTER TABLE  transaction_indexs RENAME TO transaction_indices;
COMMIT;
