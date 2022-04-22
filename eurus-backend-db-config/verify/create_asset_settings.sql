-- Verify config:create_asset_settings on pg

BEGIN;

-- XXX Add verifications here.
SELECT  *
FROM    asset_settings
LIMIT   1;

ROLLBACK;
