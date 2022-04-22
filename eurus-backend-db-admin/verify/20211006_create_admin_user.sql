-- Verify eurus-backend-db-admin:20211006_create_admin_user on pg

BEGIN;

-- XXX Add verifications here.
SELECT id FROM admin_users WHERE id = 1;

ROLLBACK;
