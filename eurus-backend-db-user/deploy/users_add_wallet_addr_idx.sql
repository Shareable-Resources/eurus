-- Deploy eurus-backend-db-user:users_add_wallet_addr_idx to pg

BEGIN;

-- XXX Add DDLs here.
CREATE UNIQUE INDEX users_wallet_addr_idx ON users (wallet_address);

COMMIT;
