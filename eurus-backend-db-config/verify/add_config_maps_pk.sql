-- Verify config:add_config_maps_pk on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT id, is_service, key FROM config_maps;
ROLLBACK;
