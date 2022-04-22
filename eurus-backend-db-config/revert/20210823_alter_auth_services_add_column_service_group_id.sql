-- Revert config:20210823_alter_auth_services_add_column_service_group_id from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE auth_services DROP COLUMN IF EXISTS service_group_id;

COMMIT;
