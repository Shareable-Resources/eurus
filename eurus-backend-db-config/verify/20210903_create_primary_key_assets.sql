-- Verify config:20210903_create_primary_key_assets on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM assets LIMIT 1;

ROLLBACK;
