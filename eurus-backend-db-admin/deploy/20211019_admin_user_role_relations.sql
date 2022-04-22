-- Deploy eurus-backend-db-admin:20211019_admin_user_role_relations to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_user_role_relations (
	admin_id bigint NOT NULL,
	role_id bigint NOT NULL,
	created_by bigint,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL,
	PRIMARY KEY (admin_id, role_id)
);
COMMIT;
