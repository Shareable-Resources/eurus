-- Verify eurus-backend-db-user:create_user_asset_setting on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM user_asset_settings LIMIT 1;

ROLLBACK;
