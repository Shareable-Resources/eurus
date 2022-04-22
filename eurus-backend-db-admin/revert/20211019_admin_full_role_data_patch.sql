-- Revert eurus-backend-db-admin:20211019_admin_full_role_data_patch from pg

BEGIN;

-- XXX Add DDLs here.

delete from admin_role_permissions where role_id in (select id from  admin_roles where role_name = 'All permission');
	
delete from admin_roles where role_name = 'All permission';

COMMIT;
