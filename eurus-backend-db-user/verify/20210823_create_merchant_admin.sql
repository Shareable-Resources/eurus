-- Verify eurus-backend-db-user:20210823_create_merchant_admin on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM merchant_admins LIMIT 1;

ROLLBACK;
