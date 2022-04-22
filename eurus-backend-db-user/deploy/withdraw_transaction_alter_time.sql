-- Deploy eurus-backend-db-user:withdraw_transaction_alter_time to pg

BEGIN;

ALTER TABLE withdraw_transactions ADD COLUMN created_date timestamptz;
ALTER TABLE withdraw_transactions ADD COLUMN last_modified_date timestamptz;

COMMIT;
