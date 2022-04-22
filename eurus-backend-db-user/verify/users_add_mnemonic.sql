-- Verify eurus-backend-db-user:users_add_mnemonic on pg

BEGIN;

-- XXX Add verifications here.
SELECT mnemonic from users limit 1;

ROLLBACK;
