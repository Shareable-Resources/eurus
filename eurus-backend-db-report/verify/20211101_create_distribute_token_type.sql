-- Verify eurus-backend-db-report:20211101_create_distribute_token_type on pg

BEGIN;

-- XXX Add verifications here.
select * from distributed_token_types limit 1;

ROLLBACK;
