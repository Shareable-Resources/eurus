-- Deploy eurus-backend-db-auth:alter_user_sessions_add_type to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_sessions ADD COLUMN IF NOT EXISTS "type" smallint ;

COMMIT;
