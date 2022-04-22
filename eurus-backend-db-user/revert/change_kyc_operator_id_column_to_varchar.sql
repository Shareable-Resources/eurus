-- Revert eurus-backend-db-user:change_kyc_operator_id_column_to_varchar from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE user_kyc_statuses
ALTER COLUMN operator_id TYPE BIGINT USING operator_id::bigint;

ALTER TABLE user_kyc_images
ALTER COLUMN operator_id TYPE BIGINT USING operator_id::bigint;

COMMIT;
