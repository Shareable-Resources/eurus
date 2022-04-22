-- Verify config:auth_services_add_column_wallet_address on pg

BEGIN;

SELECT wallet_address from auth_services LIMIT 1;

ROLLBACK;