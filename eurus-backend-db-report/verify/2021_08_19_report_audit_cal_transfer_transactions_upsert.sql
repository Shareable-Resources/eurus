-- Verify eurus-backend-db-report:2021_08_19_report_audit_cal_transfer_transactions_upsert on pg

BEGIN;

-- XXX Add verifications here.
select * from report_audit_cal_transfer_transactions_upsert('2021-08-16',2021)

ROLLBACK;
