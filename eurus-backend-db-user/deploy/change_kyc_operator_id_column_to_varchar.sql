-- Deploy eurus-backend-db-user:change_kyc_operator_id_column_to_varchar to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_kyc_statuses
ALTER COLUMN operator_id TYPE VARCHAR(25);

ALTER TABLE user_kyc_images
ALTER COLUMN operator_id TYPE VARCHAR(25);

COMMIT;
