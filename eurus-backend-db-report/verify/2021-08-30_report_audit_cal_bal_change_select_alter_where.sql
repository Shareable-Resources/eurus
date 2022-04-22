-- Verify eurus-backend-db-report:2021-08-30_report_audit_cal_bal_change_select_alter_where on pg

BEGIN;

-- XXX Add verifications here.
select * from report_audit_cal_bal_change_select('2021-08-15',90);

ROLLBACK;