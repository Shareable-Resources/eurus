-- Verify eurus-backend-db-user:20220124_alter_distributed_tokens on pg

BEGIN;

-- XXX Add verifications here.
SELECT status from distributed_tokens LIMIT 1;

ROLLBACK;
