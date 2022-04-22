-- Revert eurus-backend-db-user:transaction_indices_currency_address_to_asset_name from pg

BEGIN;


ALTER TABLE transaction_indices DROP COLUMN asset_name;
ALTER TABLE transaction_indices ADD COLUMN currency_address VARCHAR(100);

COMMIT;
