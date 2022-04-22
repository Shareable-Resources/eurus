-- Deploy eurus-backend-db-user:transfer_trans_add_gas_used to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE transfer_transactions ADD COLUMN trans_gas_used BIGINT;
ALTER TABLE transfer_transactions ADD COLUMN user_gas_used BIGINT;




COMMIT;
