-- Revert eurus-backend-db-user:create_user_kyc_images from pg
BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE IF EXISTS user_kyc_images;

COMMIT;

