-- Verify eurus-backend-db-user:create_asset_configs on pg

BEGIN;

select * from asset_configs LIMIT 1;

ROLLBACK;
