-- Deploy eurus-backend-db-user:create_admin_users to pg

BEGIN;

-- XXX Add DDLs here.

CREATE TABLE IF NOT EXISTS admin_users (
	username VARCHAR(25) NOT NULL PRIMARY KEY,
	password VARCHAR NOT NULL
);

COMMIT;
