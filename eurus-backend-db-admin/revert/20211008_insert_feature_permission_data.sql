-- Revert eurus-backend-db-admin:20211008_insert_feature_permission_data from pg

BEGIN;

-- XXX Add DDLs here.
DELETE FROM admin_feature_permissions where id >= 1 and id <= 5;

COMMIT;
