-- Verify eurus-backend-db-user:alter_users_login_address_unique on pg

BEGIN;

-- XXX Add verifications here.
SELECT login_address from users LIMIT 1;

ROLLBACK;
