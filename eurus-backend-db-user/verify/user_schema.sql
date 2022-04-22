-- Verify eurus-backend-db-user:user_schema on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM users;
ROLLBACK;
