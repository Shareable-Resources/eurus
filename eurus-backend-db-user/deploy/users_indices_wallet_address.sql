-- Deploy eurus-backend-db-user:users_indices_wallet_address to pg

BEGIN;

CREATE UNIQUE INDEX idx_users_wallet_address ON users(wallet_address);

COMMIT;
