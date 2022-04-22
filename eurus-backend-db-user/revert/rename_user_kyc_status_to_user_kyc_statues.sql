-- Revert eurus-backend-db-user:rename_user_kyc_status_to_user_kyc_statues from pg

BEGIN;

ALTER TABLE user_kyc_statuses
  RENAME TO user_kyc_status;

COMMIT;
