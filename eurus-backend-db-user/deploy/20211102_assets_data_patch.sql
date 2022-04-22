-- Deploy eurus-backend-db-user:20211102_assets_data_patch to pg

BEGIN;

-- XXX Add DDLs here.
INSERT INTO assets ("decimal", asset_name, currency_id, auto_update) VALUES (18, 'EUN', 'N/A', false);

COMMIT;
