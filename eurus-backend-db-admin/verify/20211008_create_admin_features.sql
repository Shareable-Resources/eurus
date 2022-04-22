-- Verify eurus-backend-db-admin:20211008_create_admin_features on pg

BEGIN;

-- XXX Add verifications here.
SELECT 1 FROM admin_features LIMIT 1;

ROLLBACK;
