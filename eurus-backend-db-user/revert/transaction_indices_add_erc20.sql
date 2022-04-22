-- Revert eurus-backend-db-user:transaction_indices_add_erc20 from pg

BEGIN;

ALTER TABLE transaction_indices DROP COLUMN currency_address;

COMMIT;
