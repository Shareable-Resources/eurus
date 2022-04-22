-- Revert eurus-backend-db-user:users_add_wallet_addr_idx from pg

BEGIN;

-- XXX Add DDLs here.
DROP INDEX users_wallet_addr_idx;

COMMIT;
