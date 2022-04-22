-- Verify eurus-backend-db-user:withdraw_transaction_alter_time on pg

BEGIN;

    SELECT last_modified_date,created_date FROM withdraw_transactions LIMIT 1;

ROLLBACK;
