-- Verify eurus-backend-db-admin:20211008_create_admin_feature_permission_relations on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1 FROM admin_feature_permission_relations LIMIT 1;


ROLLBACK;
