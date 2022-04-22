-- Deploy eurus-backend-db-admin:20211019_admin_full_role_data_patch to pg

BEGIN;

-- XXX Add DDLs here.
INSERT INTO admin_roles (role_name, description, state, created_date, last_modified_date) VALUES ('All permission', 'All permission', 1, NOW(), NOW());


INSERT INTO admin_role_permissions (role_id, permission_id, feature_id, created_date, last_modified_date) SELECT r.id, p.permission_id, p.feature_id, NOW(), NOW() from admin_roles as r JOIN admin_feature_permission_relations as p on 1 = 1 where r.id = ( SELECT currval('admin_roles_id_seq'));

COMMIT;
