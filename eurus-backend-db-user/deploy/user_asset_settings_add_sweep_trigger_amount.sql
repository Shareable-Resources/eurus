-- Deploy eurus-backend-db-user:user_asset_settings_add_sweep_trigger_amount to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_asset_settings ADD COLUMN IF NOT EXISTS sweep_trigger_amount NUMERIC(78);

COMMIT;
