-- Revert eurus-backend-db-report:2021-08-30_report_audit_cal_bal_change_select_alter_where from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_audit_cal_bal_change_select;


COMMIT;
