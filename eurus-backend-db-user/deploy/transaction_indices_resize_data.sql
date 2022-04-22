-- Deploy eurus-backend-db-user:transaction_indices_resize_data to pg

BEGIN;

ALTER TABLE transaction_indices ALTER COLUMN wallet_address TYPE VARCHAR(100);
ALTER TABLE transaction_indices ALTER COLUMN tx_hash TYPE VARCHAR(100);

COMMIT;
