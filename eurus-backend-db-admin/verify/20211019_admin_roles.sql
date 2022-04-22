-- Verify eurus-backend-db-admin:20211019_admin_roles on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM admin_roles LIMIT 1;

ROLLBACK;
