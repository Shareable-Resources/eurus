-- Verify eurus-backend-db-user:create_admin_users on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM admin_users

ROLLBACK;
