-- Deploy eurus-backend-db-user:withdraw_add_user_gas_used_column to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE pending_prewithdraws ADD COLUMN IF NOT EXISTS  user_gas_used BIGINT;
ALTER TABLE pending_prewithdraws ADD COLUMN IF NOT EXISTS  gas_price numeric(78);


ALTER TABLE withdraw_transactions ADD COLUMN IF NOT EXISTS  user_gas_used BIGINT;
ALTER TABLE withdraw_transactions ADD COLUMN IF NOT EXISTS  gas_price numeric(78);

COMMIT;
