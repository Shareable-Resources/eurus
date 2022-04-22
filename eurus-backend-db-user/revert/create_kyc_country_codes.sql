-- Revert eurus-backend-db-user:create_kyc_country_codes from pg
BEGIN;

-- XXX Add DDLs here.
SET search_path to public;
DROP TABLE IF EXISTS kyc_country_codes;
COMMIT;
