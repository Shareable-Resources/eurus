-- Deploy eurus-backend-db-user:create_assets to pg

BEGIN;

CREATE TABLE IF NOT EXISTS assets (
    decimal bigint,
    asset_name Varchar(100) NOT NULL
);


COMMIT;
