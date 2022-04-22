-- Revert eurus-backend-db-user:users_add_status from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users DROP COLUMN status;

COMMIT;
