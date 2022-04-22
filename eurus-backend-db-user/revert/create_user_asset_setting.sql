-- Revert eurus-backend-db-user:create_user_asset_setting from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE user_asset_settings;

COMMIT;
