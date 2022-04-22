-- Deploy eurus-backend-db-admin:20211019_admin_roles to pg

BEGIN;

-- XXX Add DDLs here.
CREATE TABLE IF NOT EXISTS admin_roles (
	id serial PRIMARY KEY,
	role_name varchar(50) NOT NULL,
	modified_by bigint,
	description varchar(255),
	state smallint NOT NULL,
	created_date timestamptz NOT NULL,
	last_modified_date timestamptz NOT NULL
);

CREATE UNIQUE INDEX admin_roles_idx1 ON admin_roles (role_name);
COMMIT;
