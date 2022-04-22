-- Verify eurus-backend-db-user:20220124_distributed_token_errors on pg

BEGIN;

-- XXX Add verifications here.
select 1 from distributed_token_errors limit 1;

ROLLBACK;
