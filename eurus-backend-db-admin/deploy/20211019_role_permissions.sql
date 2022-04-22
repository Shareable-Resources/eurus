-- Deploy eurus-backend-db-admin:20211019_role_permissions to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_role_permissions (
	role_id bigint NOT NULL,
	permission_id bigint NOT NULL,
	feature_id bigint NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL,
	PRIMARY KEY (role_id, permission_id, feature_id)
);

CREATE INDEX admin_role_permissions_idx1 ON admin_role_permissions (role_id);
COMMIT;
