-- Revert eurus-backend-db-report:2021_08_19_report_audit_cal_transfer_transactions_upsert from pg


BEGIN;

-- XXX Add DDLs here.
DROP FUNCTION IF EXISTS report_audit_cal_transfer_transactions_upsert;

COMMIT;
