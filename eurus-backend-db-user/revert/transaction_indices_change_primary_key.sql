-- Revert eurus-backend-db-user:transaction_indices_change_primary_key from pg

BEGIN;

ALTER TABLE transaction_indices DROP CONSTRAINT transaction_indexs_pkey;
ALTER TABLE transaction_indices ADD CONSTRAINT transaction_indexs_pkey primary key (user_id);

COMMIT;
