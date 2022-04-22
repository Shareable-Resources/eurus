-- Revert eurus-backend-db-user:20211102_assets_data_patch from pg

BEGIN;

-- XXX Add DDLs here.
DELETE FROM assets where asset_name = 'EUN';

COMMIT;
