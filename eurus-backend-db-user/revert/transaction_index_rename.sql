-- Revert eurus-backend-db-user:transaction_index_rename from pg

BEGIN;

ALTER TABLE transaction_indices RENAME TO transaction_indexs;

COMMIT;
