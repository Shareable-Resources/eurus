-- Revert eurus-backend-db-auth:user_sessions_alter_token from pg

BEGIN;

-- XXX Add DDLs here.
alter table user_sessions alter column token  type varchar(256);

COMMIT;
