-- Revert eurus-backend-db-user:pending_prewithdraws_add_gas_used from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE pending_prewithdraws DROP COLUMN sidechain_gas_used;


COMMIT;
