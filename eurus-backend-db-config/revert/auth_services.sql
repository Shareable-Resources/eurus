-- Revert config:auth_services from pg

BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE auth_services;

COMMIT;
