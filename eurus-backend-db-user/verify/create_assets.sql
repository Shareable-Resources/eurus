-- Verify eurus-backend-db-user:create_assets on pg

BEGIN;

select * from assets LIMIT 1;

ROLLBACK;
