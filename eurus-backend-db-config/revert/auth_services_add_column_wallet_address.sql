-- Revert config:auth_services_add_column_wallet_address from pg

BEGIN;

ALTER table auth_services DROP column wallet_address ;

COMMIT;