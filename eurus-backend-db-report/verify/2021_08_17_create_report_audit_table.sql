-- Verify eurus-backend-db-report:2021_08_17_create_report_audit_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM report_audit LIMIT 1;

ROLLBACK;
