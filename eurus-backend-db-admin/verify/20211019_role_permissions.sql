-- Verify eurus-backend-db-admin:20211019_role_permissions on pg

BEGIN;

-- XXX Add verifications here.
Select * from admin_role_permissions where 1 = 1 limit 1;

ROLLBACK;
