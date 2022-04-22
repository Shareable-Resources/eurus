-- Verify eurus-backend-db-user:pending_prewithdraws_add_admin_fee on pg

BEGIN;

SELECT admin_fee FROM pending_prewithdraws LIMIT 1;

ROLLBACK;
