-- Deploy eurus-backend-db-user:user_add_change_password_count to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users ADD column change_login_password_count INT DEFAULT 0, ADD column change_payment_password_count INT DEFAULT 0;

COMMIT;
