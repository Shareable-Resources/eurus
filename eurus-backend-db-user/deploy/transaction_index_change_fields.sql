-- Deploy eurus-backend-db-user:transaction_index_change_fields to pg

BEGIN;

ALTER TABLE transaction_indices DROP CONSTRAINT transaction_indexs_pkey;
ALTER TABLE transaction_indices ADD CONSTRAINT transaction_indices_pkey primary key (user_id, tx_hash);
COMMIT;
