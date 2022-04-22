-- Revert eurus-backend-db-user:create_assets from pg

BEGIN;

DROP TABLE assets;

COMMIT;
