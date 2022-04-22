-- Deploy eurus-backend-db-user:create_exchange_rates to pg

BEGIN;

CREATE TABLE IF NOT EXISTS exchange_rates (
    asset_name Varchar(100) PRIMARY KEY,
    rate numeric(78,18),
    created_date timestamptz,
    last_modified_date timestamptz
);

COMMIT;
