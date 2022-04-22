-- Revert eurus-backend-db-user:exchange_rates_auto_update_column_move from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE assets DROP COLUMN auto_update;
ALTER TABLE exchange_rates ADD COLUMN IF NOT EXISTS auto_update BOOL NOT NULL DEFAULT (TRUE);


COMMIT;
