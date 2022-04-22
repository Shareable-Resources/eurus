-- Deploy eurus-backend-db-admin:20211008_insert_feature_permission_data to pg

BEGIN;

-- XXX Add DDLs here.
INSERT INTO admin_feature_permissions (id, name) VALUES (1, 'Query');
INSERT INTO admin_feature_permissions (id, name) VALUES (2, 'New');
INSERT INTO admin_feature_permissions (id, name) VALUES (3, 'Update');
INSERT INTO admin_feature_permissions (id, name) VALUES (4, 'Delete');
INSERT INTO admin_feature_permissions (id, name) VALUES (5, 'Approval');

COMMIT;
