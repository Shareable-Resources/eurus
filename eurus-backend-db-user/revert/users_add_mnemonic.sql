-- Revert eurus-backend-db-user:users_add_mnemonic from pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE users DROP COLUMN mnemonic, DROP COLUMN owner_wallet_address;


COMMIT;
