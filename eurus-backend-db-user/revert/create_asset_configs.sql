-- Revert eurus-backend-db-user:create_asset_configs from pg

BEGIN;

DROP TABLE asset_configs;

COMMIT;
