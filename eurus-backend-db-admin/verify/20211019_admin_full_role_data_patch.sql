-- Verify eurus-backend-db-admin:20211019_admin_full_role_data_patch on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1/COUNT(*) FROM admin_roles where role_name = 'All permission';

SELECT 1/COUNT(*) FROM admin_roles AS a inner join admin_role_permissions as p on a.id = p.role_id where a.role_name = 'All permission';

ROLLBACK;
