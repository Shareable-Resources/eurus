-- Revert eurus-backend-db-user:withdraw_transaction_alter_time from pg

BEGIN;

    ALTER TABLE withdraw_transactions DROP COLUMN created_date,last_modified_date;

COMMIT;
