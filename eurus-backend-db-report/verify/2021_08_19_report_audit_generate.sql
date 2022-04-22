-- Verify eurus-backend-db-report:2021_08_19_report_audit_generate on pg

BEGIN;

-- XXX Add verifications here.
select report_audit_generate('2021-08-15');

ROLLBACK;
