-- Deploy config:auth_services_add_column_wallet_address to pg

BEGIN;

ALTER table auth_services ADD column wallet_address VARCHAR;
ALTER TABLE auth_services ADD column is_service bool;

COMMIT;