-- Verify eurus-backend-db-user:exchange_rate_add_auto on pg

BEGIN;

-- XXX Add verifications here.
SELECT auto_update FROM exchange_rates LIMIT 1;

ROLLBACK;
