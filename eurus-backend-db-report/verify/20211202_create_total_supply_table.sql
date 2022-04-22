-- Verify eurus-backend-db-report:20211202_create_total_supply_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM asset_total_supplies LIMIT 1;

ROLLBACK;
