-- Revert config:config_maps from pg

BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE config_maps;

COMMIT;
