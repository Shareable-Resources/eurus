-- Revert config:20211202_add_assets_eth from pg

BEGIN;

-- XXX Add DDLs here.
DELETE FROM assets where asset_name = 'ETH';

COMMIT;
