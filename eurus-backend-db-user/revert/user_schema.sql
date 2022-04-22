-- Revert eurus-backend-db-user:user_schema from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE users;
COMMIT;
