-- Revert eurus-backend-db-report:2021_08_17_create_report_audit_table from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE report_audit;

COMMIT;
