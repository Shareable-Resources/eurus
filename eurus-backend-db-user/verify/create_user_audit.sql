-- Verify eurus-backend-db-user:create_user_audit on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM user_audits;

ROLLBACK;
