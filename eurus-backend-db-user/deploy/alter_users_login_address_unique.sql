-- Deploy eurus-backend-db-user:alter_users_login_address_unique to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users ADD UNIQUE (login_address);

COMMIT;
