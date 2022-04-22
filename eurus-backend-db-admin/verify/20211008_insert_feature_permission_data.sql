-- Verify eurus-backend-db-admin:20211008_insert_feature_permission_data on pg

BEGIN;

-- XXX Add verifications here.
SELECT 5 / count(*) FROM admin_feature_permissions where id >= 1 and id <= 5;

ROLLBACK;
