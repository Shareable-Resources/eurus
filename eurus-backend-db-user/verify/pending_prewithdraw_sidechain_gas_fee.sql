-- Verify eurus-backend-db-user:pending_prewithdraw_sidechain_gas_fee on pg

BEGIN;

SELECT sidechain_gas_fee from pending_prewithdraws LIMIT 1;

ROLLBACK;
