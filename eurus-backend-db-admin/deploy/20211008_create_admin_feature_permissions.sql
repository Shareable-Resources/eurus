-- Deploy eurus-backend-db-admin:20211008_create_admin_feature_permissions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_feature_permissions (
	id bigint PRIMARY KEY,
	name varchar(255) NOT NULL
);

COMMIT;
