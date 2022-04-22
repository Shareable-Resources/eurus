-- Revert eurus-backend-db-user:users_indices_wallet_address from pg

BEGIN;

Drop INDEX idx_users_wallet_address;
COMMIT;
