-- Verify eurus-backend-db-report:2021_08_19_report_allocation_cost on pg

BEGIN;

-- XXX Add verifications here.
SELECT * from report_allocaton_cost('2021-04-19','Withdrawal')

ROLLBACK;
