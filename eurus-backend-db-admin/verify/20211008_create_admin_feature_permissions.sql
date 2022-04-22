-- Verify eurus-backend-db-admin:20211008_create_admin_feature_permissions on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1 FROM admin_feature_permissions LIMIT 1;

ROLLBACK;
