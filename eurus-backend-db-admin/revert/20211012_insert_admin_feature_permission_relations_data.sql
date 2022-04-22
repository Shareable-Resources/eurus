-- Revert eurus-backend-db-admin:20211012_insert_admin_feature_permission_relations_data from pg

BEGIN;

-- XXX Add DDLs here.
delete from admin_feature_permission_relations where feature_id >= 1 AND feature_id <= 37;

COMMIT;
