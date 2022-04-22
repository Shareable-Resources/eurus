-- Deploy eurus-backend-db-user:exchange_rates_auto_update_column_move to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE exchange_rates DROP COLUMN auto_update;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS auto_update BOOL NOT NULL DEFAULT (TRUE);


COMMIT;
