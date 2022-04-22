-- Deploy eurus-backend-db-admin:20211008_create_admin_feature_permission_relations to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_feature_permission_relations(
	feature_id bigint NOT NULL,
	permission_id bigint NOT NULL,
	PRIMARY KEY (feature_id, permission_id)
);




COMMIT;
