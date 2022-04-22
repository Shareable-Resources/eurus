-- Verify eurus-backend-db-auth:alter_user_sessions_add_type on pg

BEGIN;

-- XXX Add verifications here.
SELECT "type" from User_sessions limit 1;

ROLLBACK;
