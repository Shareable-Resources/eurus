-- Deploy config:config_maps_data_patch to pg

BEGIN;

-- XXX Add DDLs here.
UPDATE config_maps SET is_service = FALSE WHERE id = 0;

COMMIT;
