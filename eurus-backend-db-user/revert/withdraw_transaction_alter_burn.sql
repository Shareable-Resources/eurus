-- Revert eurus-backend-db-user:withdraw_transaction_alter_burn from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE withdraw_transaction DROP COLUMN burn_trans_hash,"status",burn_date,request_date;
COMMIT;
