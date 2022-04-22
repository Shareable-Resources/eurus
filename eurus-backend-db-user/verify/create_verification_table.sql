-- Verify eurus-backend-db-user:create_verification_table on pg

BEGIN;

-- XXX Add verifications here.
select * from verifications LIMIT 1;

ROLLBACK;
