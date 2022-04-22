-- Verify eurus-backend-db-user:users_alter_id on pg

BEGIN;

-- XXX Add verifications here.
SELECT ID from users limit 1;

ROLLBACK;
