-- Deploy eurus-backend-db-auth:user_sessions_alter_token to pg

BEGIN;

-- XXX Add DDLs here.
alter table user_sessions alter column token  type varchar(1024);

COMMIT;
