-- Revert eurus-backend-db-user:create_user_preference_storage from pg

BEGIN;

drop table user_preference_storages;

COMMIT;
