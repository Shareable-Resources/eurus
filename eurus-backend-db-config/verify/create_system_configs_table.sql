-- Verify config:create_system_configs_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT  *
FROM    system_configs
LIMIT   1;

ROLLBACK;
