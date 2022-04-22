-- Revert eurus-backend-db-user:transaction_index_change_fields from pg

-- Deploy eurus-backend-db-user:transaction_index_change_fields to pg

BEGIN;

ALTER TABLE transaction_indices DROP CONSTRAINT transaction_indices_pkey;
ALTER TABLE transaction_indices ADD CONSTRAINT transaction_indexs_pkey primary key (tx_hash);
COMMIT;
