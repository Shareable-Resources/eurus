-- Deploy config:20210823_alter_auth_services_add_column_service_group_id to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE auth_services ADD COLUMN service_group_id SMALLINT;

COMMIT;
