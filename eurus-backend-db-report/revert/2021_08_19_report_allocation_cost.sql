-- Revert eurus-backend-db-report:2021_08_19_report_allocation_cost from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_allocation_cost;

COMMIT;
