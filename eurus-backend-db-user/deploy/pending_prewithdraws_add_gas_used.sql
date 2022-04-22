-- Deploy eurus-backend-db-user:pending_prewithdraws_add_gas_used to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE pending_prewithdraws ADD COLUMN sidechain_gas_used numeric(78);

COMMIT;
