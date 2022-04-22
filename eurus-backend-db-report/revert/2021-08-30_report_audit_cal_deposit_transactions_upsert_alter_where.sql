-- Revert eurus-backend-db-report:2021-08-30_report_audit_cal_deposit_transactions_upsert_alter_where from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_audit_cal_deposit_transactions_upsert;

COMMIT;
