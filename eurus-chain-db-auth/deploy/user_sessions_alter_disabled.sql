-- Deploy eurus-backend-db-auth:user_sessions_alter_disabled to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_sessions ADD COLUMN disabled boolean;

COMMIT;
