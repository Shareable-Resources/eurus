-- Revert eurus-backend-db-report:2021_08_18_report_audit_cal_bal_change_select from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_audit_cal_bal_change_select;

COMMIT;
