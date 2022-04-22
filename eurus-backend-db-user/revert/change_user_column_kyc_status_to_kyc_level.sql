-- Revert eurus-backend-db-user:change_user_column_kyc_status_to_kyc_level from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users
RENAME COLUMN kyc_level TO kyc_status;
 
COMMIT;
