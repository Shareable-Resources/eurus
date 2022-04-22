-- Verify eurus-backend-db-user:user_add_change_password_count on pg

BEGIN;

-- XXX Add verifications here.
SELECT change_login_password_count, change_payment_password_count from users limit 1;

ROLLBACK;
