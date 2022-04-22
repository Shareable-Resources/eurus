-- Revert eurus-backend-db-user:pending_prewithdraw_sidechain_gas_fee from pg

BEGIN;

ALTER TABLE pending_prewithdraws DROP COLUMN sidechain_gas_fee;

COMMIT;
