-- Revert eurus-backend-db-user:alter_users_login_address_unique from pg

BEGIN;

-- XXX Add DDLs here.
alter table users drop constraint users_login_address_key;

COMMIT;
