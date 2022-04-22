-- Revert eurus-backend:user_sessions_schema from pg

BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE user_sessions;

COMMIT;
