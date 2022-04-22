-- Deploy eurus-backend-db-user:assets_add_currency_id to pg

BEGIN;

ALTER TABLE assets ADD COLUMN currency_id VARCHAR(100) NOT NULL PRIMARY KEY;

COMMIT;
