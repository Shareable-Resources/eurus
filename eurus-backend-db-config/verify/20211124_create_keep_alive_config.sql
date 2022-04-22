-- Verify config:20211124_create_keep_alive_config on pg

BEGIN;

-- XXX Add verifications here.
select * from keep_alive_config limit 1;

ROLLBACK;
