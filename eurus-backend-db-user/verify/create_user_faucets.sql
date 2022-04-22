-- Verify eurus-backend-db-user:create_user_faucets on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM user_faucets;
ROLLBACK;
