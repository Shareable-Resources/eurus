-- Revert config:create_service_groups from pg

BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE service_groups;

COMMIT;
