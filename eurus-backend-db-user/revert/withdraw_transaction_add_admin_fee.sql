-- Revert eurus-backend-db-user:withdraw_transaction_add_admin_fee from pg

BEGIN;

ALTER TABLE withdraw_transactions DROP COLUMN admin_fee;


COMMIT;
