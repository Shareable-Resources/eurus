-- Verify config:20211202_add_assets_eth on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1/count(*) FROM assets WHERE asset_name = 'ETH';

ROLLBACK;
