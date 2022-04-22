-- Verify eurus-backend-db-user:create_kyc_country_codes on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM kyc_country_codes LIMIT 1;
ROLLBACK;
