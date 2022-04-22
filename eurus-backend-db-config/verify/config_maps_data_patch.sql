-- Verify config:config_maps_data_patch on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1.0 / COUNT(*) FROM config_maps WHERE id = 0 AND is_service = FALSE;

ROLLBACK;
