-- Verify eurus-backend-db-user:20211102_assets_data_patch on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1/COUNT(*) FROM assets where asset_name = 'EUN';

ROLLBACK;
