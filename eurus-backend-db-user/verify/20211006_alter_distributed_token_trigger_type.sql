-- Verify eurus-backend-db-user:20211006_alter_distributed_token_trigger_type on pg

BEGIN;

-- XXX Add verifications here.
SELECT trigger_type from distributed_tokens LIMIT 1;

ROLLBACK;
