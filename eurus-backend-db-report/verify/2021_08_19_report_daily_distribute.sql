-- Verify eurus-backend-db-report:2021_08_19_report_daily_distribute on pg

BEGIN;

-- XXX Add verifications here.
select * from report_daily_distribute('2021-08-13',1,'EUN',2021)

ROLLBACK;
