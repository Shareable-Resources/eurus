-- Deploy eurus-backend-db-user:create_asset_configs to pg

BEGIN;

CREATE TABLE IF NOT EXISTS asset_configs (
    id serial PRIMARY KEY,
    key VARCHAR(255),
    value TEXT,
    is_service BOOLEAN
);

COMMIT;
