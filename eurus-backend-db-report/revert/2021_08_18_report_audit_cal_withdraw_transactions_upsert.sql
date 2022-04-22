-- Revert eurus-backend-db-report:2021_08_18_report_audit_cal_withdraw_transactions_upsert from pg

BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_audit_cal_withdraw_transactions_upsert;

COMMIT;
