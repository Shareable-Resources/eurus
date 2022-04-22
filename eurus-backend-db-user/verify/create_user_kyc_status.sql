-- Verify eurus-backend-db-user:create_user_kyc_status on pg

BEGIN;

SET search_path to public;
-- XXX Add verifications here.
SELECT * FROM user_kyc_statuses LIMIT 1;

ROLLBACK;
