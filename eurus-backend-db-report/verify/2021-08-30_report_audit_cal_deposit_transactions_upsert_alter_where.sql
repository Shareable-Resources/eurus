-- Verify eurus-backend-db-report:2021-08-30_report_audit_cal_deposit_transactions_upsert_alter_where on pg

BEGIN;

-- XXX Add verifications here.
select report_audit_cal_deposit_transactions_upsert('2021-08-15');


ROLLBACK;
