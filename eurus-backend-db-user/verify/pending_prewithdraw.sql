-- Verify eurus-backend-db-user:pending_prewithdraw on pg

BEGIN;

SELECT * FROM pending_prewithdraws LIMIT 1;

ROLLBACK;
