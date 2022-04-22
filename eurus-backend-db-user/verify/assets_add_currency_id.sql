-- Verify eurus-backend-db-user:assets_add_currency_id on pg

BEGIN;

SELECT currency_id from assets LIMIT 1;

ROLLBACK;
