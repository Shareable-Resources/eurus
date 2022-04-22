-- Verify eurus-backend-db-user:20210816_change_user_kyc_images_column_type_reject_reason on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM user_kyc_images LIMIT 1;

ROLLBACK;
