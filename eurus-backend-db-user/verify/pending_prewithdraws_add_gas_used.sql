-- Verify eurus-backend-db-user:pending_prewithdraws_add_gas_used on pg

BEGIN;

-- XXX Add verifications here.
SELECT sidechain_gas_used from pending_prewithdraws LIMIT 1;

ROLLBACK;
