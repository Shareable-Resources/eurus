-- Deploy eurus-backend-db-user:users_add_mnemonic to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users ADD column IF NOT EXISTS mnemonic TEXT, ADD column IF NOT EXISTS owner_wallet_address VARCHAR(50);

COMMIT;
