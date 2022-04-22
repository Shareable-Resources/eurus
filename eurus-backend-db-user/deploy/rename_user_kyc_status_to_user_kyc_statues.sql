-- Deploy eurus-backend-db-user:rename_user_kyc_status_to_user_kyc_statues to pg

BEGIN;

ALTER TABLE user_kyc_status
  RENAME TO user_kyc_statuses;

COMMIT;
