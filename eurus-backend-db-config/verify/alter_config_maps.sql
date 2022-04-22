-- Verify config:alter_config_maps on pg


SET search_path to public;
-- XXX Add verifications here.
SELECT id, is_service FROM config_maps;

ROLLBACK;

