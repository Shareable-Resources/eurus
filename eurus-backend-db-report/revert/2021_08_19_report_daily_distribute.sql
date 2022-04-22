-- Revert eurus-backend-db-report:2021_08_19_report_daily_distribute from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_daily_distribute;

COMMIT;
