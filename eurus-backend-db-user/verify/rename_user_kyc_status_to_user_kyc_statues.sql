-- Verify eurus-backend-db-user:rename_user_kyc_status_to_user_kyc_statues on pg

BEGIN;

SELECT * FROM user_kyc_statuses

ROLLBACK;
