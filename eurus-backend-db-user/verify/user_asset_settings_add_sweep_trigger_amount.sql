-- Verify eurus-backend-db-user:user_asset_settings_add_sweep_trigger_amount on pg

BEGIN;

-- XXX Add verifications here.
SELECT  *
FROM    user_asset_settings
LIMIT   1;

ROLLBACK;
