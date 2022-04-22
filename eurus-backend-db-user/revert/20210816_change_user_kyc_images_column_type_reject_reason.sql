-- Revert eurus-backend-db-user:20210816_change_user_kyc_images_column_type_reject_reason from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_kyc_images
ALTER COLUMN reject_reason TYPE VARCHAR(255);

COMMIT;
