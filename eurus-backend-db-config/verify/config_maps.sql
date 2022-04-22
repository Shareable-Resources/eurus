-- Verify config:config_maps on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM config_maps;

ROLLBACK;
