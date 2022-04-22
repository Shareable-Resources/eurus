-- Verify eurus-backend-db-user:20210804_alter_users_audits_trigger_func on pg

BEGIN;

-- XXX Add verifications here.
select kyc_level from user_audits limit 1;

ROLLBACK;
