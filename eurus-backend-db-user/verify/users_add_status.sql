-- Verify eurus-backend-db-user:users_add_status on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM users where status = 1 limit 1;

ROLLBACK;
