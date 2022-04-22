-- Verify eurus-backend-db-user:20210811_alter_pending_sweep_wallets on pg

BEGIN;

-- XXX Add verifications here.
SELECT  *
FROM    pending_sweep_wallets
LIMIT   1;

ROLLBACK;
