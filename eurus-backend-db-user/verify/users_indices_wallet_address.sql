-- Verify eurus-backend-db-user:users_indices_wallet_address on pg

BEGIN;

select * from users limit 1;

ROLLBACK;
