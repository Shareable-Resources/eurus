-- Revert eurus-backend-db-user:transaction_indices_resize_data from pg

BEGIN;

BEGIN;

ALTER TABLE transaction_indices ALTER COLUMN wallet_address TYPE VARCHAR(50);
ALTER TABLE transaction_indices ALTER COLUMN tx_hash TYPE VARCHAR(50);

COMMIT;

COMMIT;
