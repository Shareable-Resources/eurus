-- Deploy eurus-backend-db-user:pending_prewithdraw_sidechain_gas_fee to pg

BEGIN;

ALTER TABLE pending_prewithdraws ADD COLUMN sidechain_gas_fee numeric(78);

COMMIT;
