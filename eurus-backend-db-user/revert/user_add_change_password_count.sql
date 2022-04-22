-- Revert eurus-backend-db-user:user_add_change_password_count from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users DROP COLUMN change_login_password_count, DROP COLUMN change_payment_password_count;

COMMIT;
