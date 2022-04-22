-- Verify eurus-backend-db-user:20220114_alter_user_audits on pg

BEGIN;

-- XXX Add verifications here.
SELECT mnemonic from user_audits limit 1;

ROLLBACK;
