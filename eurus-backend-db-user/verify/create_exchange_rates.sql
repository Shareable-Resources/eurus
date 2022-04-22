-- Verify eurus-backend-db-user:create_exchange_rates on pg

BEGIN;

select * from exchange_rates LIMIT 1;
ROLLBACK;
