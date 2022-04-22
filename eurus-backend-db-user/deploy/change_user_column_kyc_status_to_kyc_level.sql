-- Deploy eurus-backend-db-user:change_user_column_kyc_status_to_kyc_level to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users
RENAME COLUMN kyc_status TO kyc_level;

COMMIT;
