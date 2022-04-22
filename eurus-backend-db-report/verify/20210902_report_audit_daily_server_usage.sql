-- Verify eurus-backend-db-report:20210902_report_audit_daily_server_usage on pg

BEGIN;

-- XXX Add verifications here.
select * from report_audit_daily_server_usage('20210901');

ROLLBACK;
