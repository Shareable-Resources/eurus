-- Verify eurus-backend-db-user:change_user_column_kyc_status_to_kyc_level on pg

BEGIN;

-- XXX Add verifications here.
SELECT kyc_level from users

ROLLBACK;
