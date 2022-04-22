-- Revert eurus-backend-db-user:transfer_trans_add_gas_used from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions DROP COLUMN trans_gas_used;
ALTER TABLE transfer_transactions DROP COLUMN user_gas_used;

COMMIT;
