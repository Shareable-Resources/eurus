-- Verify config:create_service_groups on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM service_groups;

ROLLBACK;

