-- Revert config:config_maps_data_patch from pg

BEGIN;

-- XXX Add DDLs here.
UPDATE config_maps SET is_service = TRUE where id = 0;

COMMIT;
