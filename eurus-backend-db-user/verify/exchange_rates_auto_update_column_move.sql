-- Verify eurus-backend-db-user:exchange_rates_auto_update_column_move on pg

BEGIN;

-- XXX Add verifications here.
SELECT auto_update from assets limit 1;

ROLLBACK;
