-- Verify eurus-backend:user_sessions_schema on pg

BEGIN;
SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM user_sessions;
ROLLBACK;
