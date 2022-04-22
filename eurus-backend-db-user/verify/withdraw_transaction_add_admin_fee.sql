-- Verify eurus-backend-db-user:withdraw_transaction_add_admin_fee on pg

BEGIN;

SELECT admin_fee from withdraw_transactions LIMIT 1;


ROLLBACK;
