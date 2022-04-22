-- Verify eurus-backend-db-user:create_user_devices_table on pg

BEGIN;

-- XXX Add verifications here.
SELECT  *
FROM    user_devices
LIMIT   1;

ROLLBACK;
