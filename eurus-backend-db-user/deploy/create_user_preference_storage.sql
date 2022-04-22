-- Deploy eurus-backend-db-user:create_user_preference_storage to pg

BEGIN;

Create table user_preference_storages (
    user_id BIGINT,
    sequence serial,
    platform SMALLINT NOT NULL,
    storage text
);

CREATE INDEX idx_user_preference_storages_userid_platform ON user_preference_storages(user_id,platform);


COMMIT;
