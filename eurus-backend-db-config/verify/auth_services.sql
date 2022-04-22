-- Verify config:auth_services on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM auth_services;

ROLLBACK;
