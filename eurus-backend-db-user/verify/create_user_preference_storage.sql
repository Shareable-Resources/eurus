-- Verify eurus-backend-db-user:create_user_preference_storage on pg

BEGIN;

select * from user_preference_storages LIMIT 1;

ROLLBACK;
