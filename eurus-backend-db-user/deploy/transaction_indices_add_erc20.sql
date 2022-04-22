-- Deploy eurus-backend-db-user:transaction_indices_add_erc20 to pg

BEGIN;

ALTER TABLE transaction_indices ADD COLUMN currency_address VARCHAR(100);

COMMIT;

