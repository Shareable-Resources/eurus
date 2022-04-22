-- Verify eurus-backend-db-user:change_kyc_operator_id_column_to_varchar on pg

BEGIN;

-- XXX Add verifications here.
SELECT operator_id from user_kyc_statuses;
SELECT operator_id from user_kyc_images;

ROLLBACK;
