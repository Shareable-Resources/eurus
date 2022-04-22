-- Revert eurus-backend-db-user:create_user_faucets from pg

BEGIN;

-- XXX Add DDLs here.
DROP TABLE user_faucets;

COMMIT;
