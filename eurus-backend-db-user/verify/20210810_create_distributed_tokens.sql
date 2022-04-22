-- Verify eurus-backend-db-user:20210810_create_distributed_tokens on pg

BEGIN;

-- XXX Add verifications here.
SELECT * FROM distributed_tokens WHERE 1 = 1 LIMIT 1;

ROLLBACK;
