-- Verify config:20210823_alter_auth_services_add_column_service_group_id on pg

BEGIN;

-- XXX Add verifications here.
SELECT service_group_id from auth_services;

ROLLBACK;
