-- Deploy eurus-backend-db-admin:20211008_create_admin_permissions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_features (
	id bigint primary key,
	name varchar(255),
	parent_feature_id bigint
);

COMMIT;
