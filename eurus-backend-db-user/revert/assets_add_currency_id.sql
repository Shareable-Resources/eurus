-- Revert eurus-backend-db-user:assets_add_currency_id from pg

BEGIN;

ALTER TABLE assets DROP COLUMN currency_id;

COMMIT;
