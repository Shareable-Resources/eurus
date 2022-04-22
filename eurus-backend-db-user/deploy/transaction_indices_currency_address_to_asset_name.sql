-- Deploy eurus-backend-db-user:transaction_indices_currency_address_to_asset_name to pg

BEGIN;

ALTER TABLE transaction_indices DROP COLUMN currency_address;
ALTER TABLE transaction_indices ADD COLUMN IF NOT EXISTS asset_name VARCHAR(50);
ALTER TABLE transaction_indices ALTER COLUMN asset_name TYPE VARCHAR(50);
COMMIT;
