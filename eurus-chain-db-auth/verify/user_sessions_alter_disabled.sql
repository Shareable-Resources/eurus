-- Verify eurus-backend-db-auth:user_sessions_alter_disabled on pg

BEGIN;

-- XXX Add verifications here.
SELECT disabled FROM user_sessions WHERE disabled IS NOT NULL;
ROLLBACK;
