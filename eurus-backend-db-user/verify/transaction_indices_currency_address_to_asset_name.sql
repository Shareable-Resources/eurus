-- Verify eurus-backend-db-user:transaction_indices_currency_address_to_asset_name on pg

BEGIN;

SELECT asset_name FROM transaction_indices limit 1;

ROLLBACK;
