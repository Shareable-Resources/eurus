-- Deploy eurus-backend-db-user:users_add_status to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users ADD column status smallint NOT NULL DEFAULT(0);

COMMIT;
