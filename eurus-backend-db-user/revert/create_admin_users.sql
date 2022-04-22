-- Revert eurus-backend-db-user:create_admin_users from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE admin_users;

COMMIT;
