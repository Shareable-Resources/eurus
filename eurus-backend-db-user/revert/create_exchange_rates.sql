-- Revert eurus-backend-db-user:create_exchange_rates from pg

BEGIN;

DROP TABLE exchange_rates;

COMMIT;
