-- Revert eurus-backend-db-auth:user_sessions_alter_disabled from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_sessions DROP COLUMN disabled;
COMMIT;
