-- Revert eurus-backend-db-user:create_users_trigger from pg

BEGIN;

-- XXX Add DDLs here.

drop TRIGGER user_trigger ON users;  

drop function update_user();

COMMIT;
