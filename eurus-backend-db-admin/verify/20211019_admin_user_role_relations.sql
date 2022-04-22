-- Verify eurus-backend-db-admin:20211019_admin_user_role_relations on pg

BEGIN;

-- XXX Add verifications here.
select * from admin_user_role_relations limit 1;

ROLLBACK;
