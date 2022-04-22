-- Deploy eurus-backend-db-user:withdraw_transaction_add_admin_fee to pg

BEGIN;

ALTER TABLE withdraw_transactions ADD COLUMN IF NOT EXISTS admin_fee numeric(78) DEFAULT 0;

COMMIT;
